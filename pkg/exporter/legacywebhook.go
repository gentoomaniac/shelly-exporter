package exporter

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

const legacyWebhookPath = "legacywebhook"

var (
	// /webhook/legacy/adu/some,other,tag/?hum=38&temp=28.00&id=shellyht-21C4B6
	pathCleanupRe = regexp.MustCompile(fmt.Sprintf("^\\/?%s\\/(.*)(\\/\\?)+.*$", legacyWebhookPath))
	// Shelly/20230913-112531/v1.14.0-gcb84623 (SHHT-1)
	userAgnetRe = regexp.MustCompile("^Shelly\\/(\\d{8}-\\d{6})\\/(v\\d+\\.\\d+\\.\\d+-\\w+) \\((.*)\\)$")
)

func (e *Exporter) legacyWebhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookCounter.Inc()

	var err error
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
		if ip == "" {
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		}
	}

	hum, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("hum")), 64)
	if err != nil {
		log.Error().Err(err).Msg("faild parsing humidity from URL")
	}
	temp, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("temp")), 64)
	if err != nil {
		log.Error().Err(err).Msg("faild parsing temperature from URL")
	}
	deviceId := strings.TrimSpace(r.URL.Query().Get("id"))

	labels := make(map[string]string)
	pathSegment := pathCleanupRe.ReplaceAll([]byte(r.URL.String()), []byte("$1"))

	for _, segment := range strings.Split(string(pathSegment), "=") {
		kv := strings.Split(segment, "/")
		if len(kv) == 2 {
			labels[kv[0]] = kv[1]
		} else {
			log.Error().Str("segment", segment).Msg("invalid key/value pair in url")
		}

	}

	constLabels := prometheus.Labels{
		"userAgent": r.Header.Get("user-agent"),
		"ip":        ip,
	}
	for k, v := range labels {
		constLabels[k] = v
	}

	deviceType := string(userAgnetRe.ReplaceAll([]byte(r.Header.Get("User-Agent")), []byte("$3")))

	err = e.updateCollector(deviceType, deviceId, "temperature", temp, labels)
	if err != nil {
		log.Error().Err(err).Msg("failed registering new metric")
		log.Error().Str("remoteAddr", r.RemoteAddr).Str("sourceIp", r.Header.Get("X-Forwarded-For")).Str("uri", r.RequestURI).Str("userAgent", r.Header.Get("User-Agent")).Msg("")
	}
	err = e.updateCollector(deviceType, deviceId, "humidity", hum, labels)
	if err != nil {
		log.Error().Err(err).Msg("failed registering new metric")
		log.Error().Str("remoteAddr", r.RemoteAddr).Str("sourceIp", r.Header.Get("X-Forwarded-For")).Str("uri", r.RequestURI).Str("userAgent", r.Header.Get("User-Agent")).Msg("")
	}
}

func (e *Exporter) updateCollector(deviceType string, deviceId string, metric string, value float64, labels prometheus.Labels) error {
	monName := strings.ToLower(fmt.Sprintf("%s_%s_%s", deviceType, deviceId, metric))
	mon, ok := e.webhookCollectors[monName]
	if !ok {
		e.webhookCollectors[monName] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   "shelly",
			Name:        metric,
			ConstLabels: labels,
		},
		)
		mon = e.webhookCollectors[monName]

		err := prometheus.Register(mon)
		if err != nil {
			return err
		}
		log.Debug().Str("id", deviceId).Str("type", deviceType).Str("monitorName", monName).Msg("new collector registered")
	}
	mon.Set(value)
	log.Debug().Str("id", deviceId).Str("type", deviceType).Msg("value updated")

	return nil
}
