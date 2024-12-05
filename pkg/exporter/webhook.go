package exporter

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

func (e *Exporter) webhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookCounter.Inc()

	monitorName := strings.ToLower(fmt.Sprintf("%s_%s_%s", r.URL.Query().Get("type"), r.URL.Query().Get("deviceid"), r.URL.Query().Get("metric")))
	value, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("value")), 64)
	if err != nil {
		log.Error().Err(err).Msg("faild parsing value from URL")
	}

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

	constLabels := prometheus.Labels{
		"type":      r.URL.Query().Get("type"),
		"deviceID":  r.URL.Query().Get("deviceid"),
		"name":      r.URL.Query().Get("name"),
		"userAgent": r.Header.Get("user-agent"),
		"ip":        ip,
	}

	mon, ok := e.webhookCollectors[monitorName]
	if !ok {
		e.webhookCollectors["power_current"] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   r.URL.Query().Get("namespace"),
			Name:        r.URL.Query().Get("metric"),
			ConstLabels: constLabels,
		},
		)
		mon = e.webhookCollectors["power_current"]

		err := prometheus.Register(mon)
		if err != nil {
			log.Error().Msg("failed registering new metric")
			log.Error().Str("remoteAddr", r.RemoteAddr).Str("sourceIp", r.Header.Get("X-Forwarded-For")).Str("uri", r.RequestURI).Str("userAgent", r.Header.Get("User-Agent")).Msg("")
		}
	}

	mon.Set(value)
}
