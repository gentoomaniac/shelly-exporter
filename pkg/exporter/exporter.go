package exporter

import (
	"fmt"
	"net/http"
	"net/netip"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gentoomaniac/shelly-exporter/pkg/config"
	homewizard_v1 "github.com/gentoomaniac/shelly-exporter/pkg/homewizard/v1"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
	shelly_auth "github.com/gentoomaniac/shelly-exporter/pkg/shelly/auth"
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

			err := d.RefreshDeviceinfo()
			if err != nil {
				log.Error().Err(err).Str("device", d.Hostname()).Msg("deviceinfo refresh failed")
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
	log.Debug().Msg("finished setting up devices")

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
		case config.SHELLY:
			ip := netip.MustParseAddr(dev.IP.String())
			exporterDev, err = shelly.DeviceFromIP(&ip, &shelly_auth.Auth{User: string(user), Password: string(password)}, dev.Labels)

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
