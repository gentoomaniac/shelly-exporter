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

var requiredFields = []string{"type", "deviceid", "metric", "value"}

func (e *Exporter) webhookHandler(w http.ResponseWriter, r *http.Request) {
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

	for _, fieldName := range requiredFields {
		if r.URL.Query().Get(fieldName) == "" {
			log.Debug().Str("URI", r.RequestURI).Str("source", ip).Msgf("required field `%s` not in query parameters", fieldName)
			w.Write([]byte(fmt.Sprintf("required field `%s` not in query parameters", fieldName)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	monitorName := strings.ToLower(fmt.Sprintf("%s_%s_%s", r.URL.Query().Get("type"), r.URL.Query().Get("deviceid"), r.URL.Query().Get("metric")))
	value, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("value")), 64)
	if err != nil {
		log.Error().Err(err).Msg("faild parsing value from URL")
	}

	constLabels := prometheus.Labels{
		"userAgent": r.Header.Get("user-agent"),
		"ip":        ip,
	}
	for k, v := range r.URL.Query() {
		constLabels[k] = strings.Join(v, ",")
	}
	log.Debug().Str("labels", fmt.Sprintf("%v", constLabels)).Msg("")

	mon, ok := e.webhookCollectors[monitorName]
	if !ok {
		e.webhookCollectors[monitorName] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   r.URL.Query().Get("namespace"),
			Name:        r.URL.Query().Get("metric"),
			ConstLabels: constLabels,
		},
		)
		mon = e.webhookCollectors[monitorName]

		err := prometheus.Register(mon)
		if err != nil {
			log.Error().Err(err).Msg("failed registering new metric")
			log.Error().Str("remoteAddr", r.RemoteAddr).Str("sourceIp", r.Header.Get("X-Forwarded-For")).Str("uri", r.RequestURI).Str("userAgent", r.Header.Get("User-Agent")).Msg("")
		}
	}

	mon.Set(value)
}
