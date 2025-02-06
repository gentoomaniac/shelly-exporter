package minipmg3

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gentoomaniac/shelly-exporter/pkg/collector"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/minipmg3/api"
)

const TypeString = "SHMINIPMG3"

type Config struct {
	BaseUrl *url.URL
	Labels  map[string]string
	Ip      net.IP
	Auth    shelly.Auth
}

func NewMiniPMG3(c Config) (*MiniPMG3, error) {
	if c.BaseUrl == nil {
		if c.Ip == nil {
			return nil, fmt.Errorf("must provide at least one of `Baseurl` or `IP`")
		}
		c.BaseUrl, _ = url.Parse(fmt.Sprintf("http://%s/rpc/", c.Ip.String()))
	}

	p := &MiniPMG3{
		config:     c,
		collectors: make(map[string]prometheus.Collector),
	}

	err := p.RefreshDeviceinfo()
	if err != nil {
		return nil, err
	}

	return p, nil
}

type MiniPMG3 struct {
	config Config

	configData api.Config
	statusData api.Status

	collectors map[string]prometheus.Collector
}

func (m MiniPMG3) Name() string {
	return m.configData.Sys.Device.Name
}

func (m MiniPMG3) Hostname() string {
	return fmt.Sprintf("shellyminipmg3-%s", strings.ToLower(m.configData.Sys.Device.Mac))
}

func (m *MiniPMG3) RefreshDeviceinfo() error {
	settingsUrl := m.config.BaseUrl.JoinPath("Shelly.GetConfig")
	resp, err := shelly.DigestAuthedRequest(settingsUrl, &m.config.Auth, map[string]string{"id": "0"})
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &m.configData)
	if err != nil {
		return err
	}

	return nil
}

func (m *MiniPMG3) Refresh() error {
	statusUrl := m.config.BaseUrl.JoinPath("Shelly.GetStatus")
	resp, err := shelly.DigestAuthedRequest(statusUrl, &m.config.Auth, map[string]string{"id": "0"})
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &m.statusData)
	if err != nil {
		return err
	}

	return nil
}

func (m *MiniPMG3) Collectors() ([]prometheus.Collector, error) {
	bool2int := map[bool]int8{false: 0, true: 1}

	constLabels := prometheus.Labels{
		"type":   TypeString,
		"serial": m.configData.Sys.Device.Mac,
	}
	dynamicLabels := []string{"name", "hostname"}

	for k, v := range m.config.Labels {
		constLabels[k] = v
	}

	// Power
	m.collectors["power_current"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "power_current",
		Help:          "Current real AC power being drawn, [W]",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Pm10.Apower) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["total_energy"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "total_energy",
		Help:          "Last counter value of the total energy consumed in Watt-hours",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Pm10.Aenergy.Total) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	// System
	m.collectors["uptime"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "uptime",
		Help:          "device uptime",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.Uptime) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["memory_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "memory_total",
		Help:          "total device memory",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.RAMSize) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["memory_free"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "memory_free",
		Help:          "free device memory",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.RAMFree) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["fs_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "fs_total",
		Help:          "total filesystem size",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.FsSize) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["fs_free"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "fs_free",
		Help:          "free filesystem size",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.FsFree) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["has_update"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "has_update",
		Help:          "device update available",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(bool2int[m.statusData.Sys.AvailableUpdates.Stable.Version != ""]) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	m.collectors["has_beta_update"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "has_beta_update",
		Help:          "device beta version available",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(bool2int[m.statusData.Sys.AvailableUpdates.Beta.Version != ""]) },
		func() []string { return []string{m.configData.Sys.Device.Name, m.Hostname()} },
	)

	var c []prometheus.Collector
	for _, v := range m.collectors {
		c = append(c, v)
	}

	return c, nil
}
