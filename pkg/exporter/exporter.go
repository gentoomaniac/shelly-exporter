package exporter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gentoomaniac/shelly-exporter/pkg/config"
	homewizard_v1 "github.com/gentoomaniac/shelly-exporter/pkg/homewizard/v1"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
	shelly_minipm3g "github.com/gentoomaniac/shelly-exporter/pkg/shelly/minipmg3"
	shelly_plugs "github.com/gentoomaniac/shelly-exporter/pkg/shelly/plugs"
	shelly_pro3em "github.com/gentoomaniac/shelly-exporter/pkg/shelly/pro3em"
)

const (
	metadataRefreshInterval   = time.Minute * 5
	sensordataRefreshInterval = time.Second * 5
)

var webhookCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "shellyexporter",
	Name:      "webhook_calls",
	Help:      "Number calls to the webhook",
},
)

type Device interface {
	Collectors() ([]prometheus.Collector, error)
	Hostname() string
	Name() string
	Refresh() error
	RefreshDeviceinfo() error
}

type Exporter struct {
	config     *config.Config
	devices    []Device
	collectors map[string][]prometheus.Collector

	webhookCollectors map[string]prometheus.Gauge
}

func New(c *config.Config) *Exporter {
	return &Exporter{
		config:            c,
		collectors:        make(map[string][]prometheus.Collector),
		webhookCollectors: make(map[string]prometheus.Gauge),
	}
}

func (e *Exporter) updateDevice(d Device) {
	lastMetadataUpdate := time.Now().UTC()
	for {
		err := d.Refresh()
		if err != nil {
			log.Error().Err(err).Str("device", d.Hostname()).Msg("refresh failed")
		}
		log.Debug().Str("device", d.Hostname()).Msg("refreshed")

		if time.Now().UTC().Sub(lastMetadataUpdate) > metadataRefreshInterval {
			oldName := d.Name()
			oldHostname := d.Hostname()

			err := d.RefreshDeviceinfo()
			if err != nil {
				log.Error().Err(err).Str("device", d.Hostname()).Msg("deviceinfo refresh failed")
			}

			// TODO: This doesn't work and leaves behind an orphaned metric in the exporter
			// for reference:
			// https://stackoverflow.com/a/77900920
			if d.Name() != oldName || d.Hostname() != oldHostname {
				collectors, err := d.Collectors()
				if err != nil {
					log.Error().Err(err).Msg("failed registering collectors")
				}

				for _, c := range e.collectors[d.Name()] {
					prometheus.Unregister(c)
				}
				prometheus.MustRegister(collectors...)
				log.Debug().Str("device", d.Hostname()).Msg("collectors refreshed")
			}
			lastMetadataUpdate = time.Now().UTC()
			log.Debug().Str("device", d.Hostname()).Msg("deviceinfo refreshed")
		}

		time.Sleep(sensordataRefreshInterval)
	}
}

func (e *Exporter) Run() {
	if err := e.setupDevices(); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	var err error
	for _, dev := range e.devices {
		e.collectors[dev.Name()], err = dev.Collectors()
		if err != nil {
			log.Error().Err(err).Msg("failed registering collectors")
		}

		prometheus.MustRegister(e.collectors[dev.Name()]...)
		go e.updateDevice(dev)
	}
	prometheus.MustRegister(webhookCounter)

	// Start HTTP server
	log.Info().Msg("Starting server on port 8080")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/webhook", e.webhookHandler)
	http.HandleFunc("/legacywebhook/", e.legacyWebhookHandler)
	http.ListenAndServe(":8080", nil)
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

		if dev.Labels == nil {
			dev.Labels = make(map[string]string)
		}
		dev.Labels["ip"] = dev.IP.String()

		switch dev.Type {
		case config.SHPLG_S:
			exporterDev, err = shelly_plugs.NewPlugS(dev.IP, string(user), string(password), dev.Labels)

		case config.SHMINIPMG3:
			exporterDev, err = shelly_minipm3g.NewMiniPMG3(
				shelly_minipm3g.Config{Ip: dev.IP, Auth: shelly.Auth{User: string(user), Password: string(password)}, Labels: dev.Labels},
			)

		case config.SHPRO3EM:
			exporterDev, err = shelly_pro3em.NewPro3EM(
				shelly_pro3em.Config{Ip: dev.IP, Auth: shelly.Auth{User: string(user), Password: string(password)}, Labels: dev.Labels},
			)

		case config.HWE_P1:
			exporterDev, err = homewizard_v1.NewP1(
				homewizard_v1.Config{Ip: dev.IP, Labels: dev.Labels},
			)

		default:
			return fmt.Errorf("unknown device: %s", dev.Type)
		}

		if err != nil {
			log.Error().Err(err).Str("ip", dev.IP.String()).Msgf("failed creating device")
		} else {
			e.devices = append(e.devices, exporterDev)
		}
	}

	return nil
}
