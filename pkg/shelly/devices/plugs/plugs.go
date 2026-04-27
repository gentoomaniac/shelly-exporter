package plugs

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/netip"
	"net/url"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gentoomaniac/shelly-exporter/pkg/collector"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/auth"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/devices/plugs/api"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/request"
)

const TypeString = "SHPLG-S"

type Config struct {
	BaseUrl *url.URL
	Labels  map[string]string
	Ip      *netip.Addr
	Auth    *auth.Auth
}

func NewPlugS(c Config) (*PlugS, error) {
	if c.BaseUrl == nil {
		if c.Ip == nil {
			return nil, fmt.Errorf("must provide at least one of `Baseurl` or `IP`")
		}
		c.BaseUrl, _ = url.Parse(fmt.Sprintf("http://%s/", c.Ip.String()))
	}

	p := &PlugS{
		config:     c,
		collectors: make(map[string]prometheus.Collector),
	}

	err := p.RefreshDeviceinfo()
	if err != nil {
		return nil, err
	}

	return p, nil
}

type PlugS struct {
	config Config

	settings api.Settings
	status   api.Status

	collectors map[string]prometheus.Collector
}

func (p PlugS) Name() string {
	return p.settings.Name
}

func (p PlugS) Hostname() string {
	return p.settings.Device.Hostname
}

func (p *PlugS) RefreshDeviceinfo() error {
	settingsUrl := p.config.BaseUrl.JoinPath("settings")
	resp, err := request.Request(settingsUrl, p.config.Auth)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &p.settings)
	if err != nil {
		return err
	}

	return nil
}

func (p *PlugS) Refresh() error {
	statusUrl := p.config.BaseUrl.JoinPath("status")
	resp, err := request.Request(statusUrl, p.config.Auth)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &p.status)
	if err != nil {
		return err
	}

	return nil
}

func (p *PlugS) Collectors() ([]prometheus.Collector, error) {
	bool2int := map[bool]int8{false: 0, true: 1}

	constLabels := prometheus.Labels{
		"type":   TypeString,
		"serial": strconv.Itoa(p.status.Serial),
	}
	dynamicLabels := []string{"name", "hostname"}

	maps.Copy(constLabels, p.config.Labels)

	// Power
	p.collectors["power_current"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "power_current",
		Help:          "Current real AC power being drawn, [W]",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.Meters[0].Power) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["power_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "power_total",
		Help:          "Total energy consumed by the attached electrical appliance in Watt-minute",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.Meters[0].Total) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	// Temperatures
	p.collectors["temperature_celsius"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "temperature_celsius",
		Help:          "internal device temperature in °C",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.Tmp.Celsius) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["temperature_fahrenheit"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "temperature_fahrenheit",
		Help:          "internal device temperature in °F",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.Tmp.Fahrenheit) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	// System
	p.collectors["uptime"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "uptime",
		Help:          "device uptime",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.Uptime) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["memory_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "memory_total",
		Help:          "total device memory",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.RamTotal) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["memory_free"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "memory_free",
		Help:          "free device memory",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.RamFree) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["fs_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "fs_total",
		Help:          "total filesystem size",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.FsSize) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["fs_free"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "fs_free",
		Help:          "free filesystem size",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.status.FsFree) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	p.collectors["has_update"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "has_update",
		Help:          "device update available",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(bool2int[p.status.HasUpdate]) },
		func() []string { return []string{p.settings.Name, p.settings.Device.Hostname} },
	)

	var c []prometheus.Collector
	for _, v := range p.collectors {
		c = append(c, v)
	}

	return c, nil
}
