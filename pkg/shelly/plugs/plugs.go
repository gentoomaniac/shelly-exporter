package plugs

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gentoomaniac/shelly-exporter/pkg/collector"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
)

const TypeString = "SHPLG-S"

func NewPlugS(ip net.IP, user string, password string, labels prometheus.Labels) (*PlugS, error) {
	baseUrl, _ := url.Parse(fmt.Sprintf("http://%s/", ip.String()))

	p := &PlugS{
		baseUrl:    baseUrl,
		auth:       shelly.Auth{User: user, Password: password},
		labels:     labels,
		collectors: make(map[string]prometheus.Collector),
	}

	err := p.RefreshDeviceinfo()
	if err != nil {
		return nil, err
	}

	return p, nil
}

type PlugS struct {
	baseUrl *url.URL
	auth    shelly.Auth
	labels  prometheus.Labels

	settings shelly.Settings
	status   Status

	collectors map[string]prometheus.Collector
}

type Status struct {
	ActionStats       shelly.ActionStats `json:"action_stats"`
	ConfigChangeCount int                `json:"cfg_change_count"`
	Cloud             shelly.Cloud       `json:"cloud"`
	FsFree            int                `json:"fs_free"`
	FsSize            int                `json:"fs_size"`
	HasUpdate         bool               `json:"has_update"`
	Mac               string             `json:"mac"`
	Meters            []shelly.Meter     `json:"meters"`
	Mqtt              shelly.Mqtt        `json:"mqtt"`
	OverTemperature   bool               `json:"overtemperature"`
	RamFree           int                `json:"ram_free"`
	RamTotal          int                `json:"ram_total"`
	Relays            []shelly.Relay     `json:"relays"`
	Serial            int                `json:"serial"`
	Temperature       float32            `json:"temperature"`
	Tmp               shelly.Temperature `json:"tmp"`
	Time              string             `json:"time"`
	Unixtime          int                `json:"unixtime"`
	Update            shelly.Update      `json:"update"`
	Uptime            int                `json:"uptime"`
	Wifi              shelly.Wifi        `json:"wifi_sta"`
}

func (p PlugS) Name() string {
	return p.settings.Name
}

func (p PlugS) Hostname() string {
	return p.settings.Device.Hostname
}

func (p *PlugS) RefreshDeviceinfo() error {
	settingsUrl := p.baseUrl.JoinPath("settings")
	resp, err := shelly.Request(settingsUrl, &p.auth)
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
	statusUrl := p.baseUrl.JoinPath("status")
	resp, err := shelly.Request(statusUrl, &p.auth)
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

	for k, v := range p.labels {
		constLabels[k] = v
	}

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
