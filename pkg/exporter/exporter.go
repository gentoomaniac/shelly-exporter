package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gentoomaniac/shelly-exporter/pkg/config"
	homewizard_v1 "github.com/gentoomaniac/shelly-exporter/pkg/homewizard/v1"
	plugs "github.com/gentoomaniac/shelly-exporter/pkg/shelly/plugs"
)

var webhookCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "shellyexporter",
	Name:      "webhook_calls",
	Help:      "Number calls to the webhook",
},
)

type Device interface {
	Collectors() ([]prometheus.Collector, error)
	Name() string
	Refresh() error
	RefreshDeviceinfo() error
}

type Exporter struct {
	config  *config.Config
	devices []Device
}

func New(c *config.Config) *Exporter {
	return &Exporter{
		config: c,
	}
}

func updateDevice(s Device) {
	for {
		err := s.Refresh()
		if err != nil {
			log.Error().Err(err).Str("device", s.Name()).Msg("refresh failed")
		}
		log.Debug().Str("device", s.Name()).Msg("refreshed")
		time.Sleep(5 * time.Second)
	}
}

func (e *Exporter) Run() {
	if err := e.setupDevices(); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	for _, dev := range e.devices {
		collectors, err := dev.Collectors()
		if err != nil {
			log.Error().Err(err).Msg("failed registering collectors")
		}

		prometheus.MustRegister(collectors...)
		go updateDevice(dev)
	}
	prometheus.MustRegister(webhookCounter)

	// Start HTTP server
	log.Info().Msg("Starting server on port 8080")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/webhook", e.webhookHandler)
	http.ListenAndServe(":8080", nil)
}

func (e *Exporter) webhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookCounter.Inc()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed reading body")
	}

	var body interface{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Info().Str("remoteAddr", r.RemoteAddr).Str("uri", r.RequestURI).Any("headers", r.Header).Str("body", string(bodyBytes)).Msg("")
	} else {
		log.Info().Str("remoteAddr", r.RemoteAddr).Str("uri", r.RequestURI).Any("headers", r.Header).Any("body", body).Msg("")
	}
}

func (e *Exporter) setupDevices() (err error) {
	for _, dev := range e.config.Devices {
		var exporterDev Device

		user := e.config.Global.User
		if dev.User != "" {
			user = dev.User
		}

		password := e.config.Global.Password
		if dev.Password != "" {
			password = dev.Password
		}

		dev.Labels["ip"] = dev.IP.String()

		switch dev.Type {
		case config.SHPLG_S:
			exporterDev, err = plugs.NewPlugS(dev.IP, string(user), string(password), dev.Labels)

		case config.HWE_P1:
			exporterDev, err = homewizard_v1.NewP1(
				homewizard_v1.Config{IP: dev.IP, Labels: dev.Labels},
			)

		default:
			return fmt.Errorf("unknown device: %s", dev.Type)
		}

		if err != nil {
			return fmt.Errorf("failed creating device for ip %s: %w", dev.IP, err)
		}

		e.devices = append(e.devices, exporterDev)
	}

	return nil
}
