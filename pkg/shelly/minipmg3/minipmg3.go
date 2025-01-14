package minipmg3

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

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
		"type":     TypeString,
		"serial":   m.configData.Sys.Device.Mac,
		"name":     m.configData.Sys.Device.Name,
		"hostname": fmt.Sprintf("shellyminipmg3-%s", strings.ToLower(m.configData.Sys.Device.Mac)),
	}

	for k, v := range m.config.Labels {
		constLabels[k] = v
	}

	// Power
	m.collectors["power_current"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "power_current",
		Help:        "Current real AC power being drawn, [W]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Pm10.Apower) },
	)

	m.collectors["total_energy"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "total_energy",
		Help:        "Last counter value of the total energy consumed in Watt-hours",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Pm10.Apower) },
	)

	// System
	m.collectors["uptime"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "uptime",
		Help:        "device uptime",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.Uptime) },
	)

	m.collectors["memory_total"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "memory_total",
		Help:        "total device memory",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.RAMSize) },
	)

	m.collectors["memory_free"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "memory_free",
		Help:        "free device memory",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.RAMFree) },
	)

	m.collectors["fs_total"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "fs_total",
		Help:        "total filesystem size",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.FsSize) },
	)

	m.collectors["fs_free"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "fs_free",
		Help:        "free filesystem size",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(m.statusData.Sys.FsFree) },
	)

	m.collectors["has_update"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "has_update",
		Help:        "device update available",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(bool2int[m.statusData.Sys.AvailableUpdates.Stable.Version != ""]) },
	)

	m.collectors["has_beta_update"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "has_beta_update",
		Help:        "device beta version available",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(bool2int[m.statusData.Sys.AvailableUpdates.Beta.Version != ""]) },
	)

	var c []prometheus.Collector
	for _, v := range m.collectors {
		c = append(c, v)
	}

	return c, nil
}
