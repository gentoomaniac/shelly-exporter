package exporter

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
)

type ShellyDevice interface {
	Name() string
	Refresh() error
	Collectors() ([]prometheus.Collector, error)
}

type Exporter struct {
	config  *Config
	devices []ShellyDevice
}

func New(c *Config) *Exporter {
	return &Exporter{
		config: c,
	}
}

func updateDevice(s ShellyDevice) {
	for true {
		s.Refresh()
		log.Debug().Str("device", "").Msg("updated device")
		time.Sleep(5 * time.Second)
	}
}

func (e *Exporter) Run() {
	e.setupDevices()

	for _, dev := range e.devices {
		collectors, err := dev.Collectors()
		if err != nil {
			log.Error().Err(err).Msg("failed registering collectors")
		}

		prometheus.MustRegister(collectors...)
		go updateDevice(dev)
	}

	// Start HTTP server
	log.Info().Msg("Starting server on port 8080")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func (e *Exporter) setupDevices() {
	for _, dev := range e.config.Devices {
		var shellyDev ShellyDevice

		user := e.config.Global.User
		if dev.User != "" {
			user = dev.User
		}

		password := e.config.Global.Password
		if dev.Password != "" {
			password = dev.Password
		}

		switch dev.Type {
		case SHPLG_S:
			shellyDev = shelly.NewPlugS(dev.Ip, string(user), string(password), dev.Labels)
		}

		e.devices = append(e.devices, shellyDev)
		err := shellyDev.Refresh()
		if err != nil {
			log.Error().Err(err).Str("device", shellyDev.Name()).Msg("refresh failed")
		}
	}
}
