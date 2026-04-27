package outdoorplugsg3

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/netip"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gentoomaniac/shelly-exporter/pkg/collector"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/auth"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/devices/outdoorplugsg3/api"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/request"
)

const TypeString = "OutdoorPlugSG3"

type Config struct {
	BaseUrl *url.URL
	Labels  map[string]string
	Ip      *netip.Addr
	Auth    *auth.Auth
}

func NewOutdoorPlugSG3(c Config) (*OutdoorPlugSG3, error) {
	if c.BaseUrl == nil {
		if c.Ip == nil {
			return nil, fmt.Errorf("must provide at least one of `Baseurl` or `IP`")
		}
		c.BaseUrl, _ = url.Parse(fmt.Sprintf("http://%s/rpc/", c.Ip.String()))
	}

	p := &OutdoorPlugSG3{
		config:     c,
		collectors: make(map[string]prometheus.Collector),
	}

	err := p.RefreshDeviceinfo()
	if err != nil {
		return nil, err
	}

	return p, nil
}

type OutdoorPlugSG3 struct {
	config Config

	configData api.Config
	statusData api.Status

	collectors map[string]prometheus.Collector
}

func (p OutdoorPlugSG3) Name() string {
	if p.configData.Switch0.Name != nil {
		return *p.configData.Switch0.Name
	}
	return ""
}

func (p OutdoorPlugSG3) Hostname() string {
	return fmt.Sprintf("shellyoutdoorsg3-%s", strings.ToLower(p.configData.Sys.Device.Mac))
}

func (p *OutdoorPlugSG3) RefreshDeviceinfo() error {
	settingsUrl := p.config.BaseUrl.JoinPath("Shelly.GetConfig")
	resp, err := request.DigestAuthedRequest(settingsUrl, p.config.Auth, map[string]string{"id": "0"})
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &p.configData)
	if err != nil {
		return err
	}

	return nil
}

func (p *OutdoorPlugSG3) Refresh() error {
	statusUrl := p.config.BaseUrl.JoinPath("Shelly.GetStatus")
	resp, err := request.DigestAuthedRequest(statusUrl, p.config.Auth, map[string]string{"id": "0"})
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &p.statusData)
	if err != nil {
		return err
	}

	return nil
}

func (p *OutdoorPlugSG3) Collectors() ([]prometheus.Collector, error) {
	bool2int := map[bool]int8{false: 0, true: 1}

	constLabels := prometheus.Labels{
		"type":   TypeString,
		"serial": p.configData.Sys.Device.Mac,
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
		func() float64 { return float64(p.statusData.Switch0.APower) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["total_energy"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "total_energy",
		Help:          "Last counter value of the total energy consumed in Watt-hours",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Switch0.AEnergy.Total) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)
	// for compatibility in combined graphing with PlugS and others
	p.collectors["power_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "power_total",
		Help:          "Total energy consumed by the attached electrical appliance in Watt-minute",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Switch0.AEnergy.Total * 60) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	// System
	p.collectors["uptime"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "uptime",
		Help:          "device uptime",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Sys.Uptime) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["memory_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "memory_total",
		Help:          "total device memory",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Sys.RamSize) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["memory_free"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "memory_free",
		Help:          "free device memory",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Sys.RamFree) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["fs_total"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "fs_total",
		Help:          "total filesystem size",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Sys.FsSize) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["fs_free"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "fs_free",
		Help:          "free filesystem size",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(p.statusData.Sys.FsFree) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["has_update"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "has_update",
		Help:          "device update available",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(bool2int[p.statusData.Sys.AvailableUpdates.Stable.Version != ""]) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	p.collectors["has_beta_update"] = collector.NewDynamicLabelGaugeCollector(collector.DynamicLabelGaugeCollectorOpts{
		Namespace:     "shelly",
		Name:          "has_beta_update",
		Help:          "device beta version available",
		DynamicLabels: dynamicLabels,
		ConstLabels:   constLabels,
	},
		func() float64 { return float64(bool2int[p.statusData.Sys.AvailableUpdates.Beta.Version != ""]) },
		func() []string { return []string{p.Name(), p.Hostname()} },
	)

	var c []prometheus.Collector
	for _, v := range p.collectors {
		c = append(c, v)
	}

	return c, nil
}
