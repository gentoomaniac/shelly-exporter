package pro3em

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gentoomaniac/shelly-exporter/pkg/shelly"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/pro3em/api"
)

const TypeString = "SHPRO3EM"

type Config struct {
	BaseUrl *url.URL
	Labels  map[string]string
	Ip      net.IP
	Auth    shelly.Auth
}

func NewPro3EM(c Config) (*Pro3EM, error) {
	if c.BaseUrl == nil {
		if c.Ip == nil {
			return nil, fmt.Errorf("must provide at least one of `Baseurl` or `IP`")
		}
		c.BaseUrl, _ = url.Parse(fmt.Sprintf("http://%s/rpc/", c.Ip.String()))
	}

	p := &Pro3EM{
		config:     c,
		collectors: make(map[string]prometheus.Collector),
	}

	err := p.RefreshDeviceinfo()
	if err != nil {
		return nil, err
	}

	return p, nil
}

type Pro3EM struct {
	config Config

	configData api.Config
	statusData api.Status

	collectors map[string]prometheus.Collector
}

func (p Pro3EM) Name() string {
	return p.configData.Sys.Device.Name
}

func (p Pro3EM) Hostname() string {
	return fmt.Sprintf("shellypro3em-%s", strings.ToLower(p.configData.Sys.Device.Mac))
}

func (p *Pro3EM) RefreshDeviceinfo() error {
	settingsUrl := p.config.BaseUrl.JoinPath("Shelly.GetConfig")
	resp, err := shelly.DigestAuthedRequest(settingsUrl, &p.config.Auth, map[string]string{"id": "0"})
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &p.configData)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pro3EM) Refresh() error {
	statusUrl := p.config.BaseUrl.JoinPath("Shelly.GetStatus")
	resp, err := shelly.DigestAuthedRequest(statusUrl, &p.config.Auth, map[string]string{"id": "0"})
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &p.statusData)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pro3EM) Collectors() ([]prometheus.Collector, error) {
	// bool2int := map[bool]int8{false: 0, true: 1}

	constLabels := prometheus.Labels{
		"type":     TypeString,
		"serial":   p.configData.Sys.Device.Mac,
		"name":     p.configData.Sys.Device.Name,
		"hostname": p.Hostname(),
	}

	for k, v := range p.config.Labels {
		constLabels[k] = v
	}

	// Power
	p.collectors["total_active_power"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "total_active_power",
		Help:        "Sum of the active power on all phases, [W]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(p.statusData.Em0.TotalActPower) },
	)
	p.collectors["a_act_power"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "a_act_power",
		Help:        "Phase A active power measurement value, [W]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(p.statusData.Em0.AActPower) },
	)
	p.collectors["b_act_power"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "b_act_power",
		Help:        "Phase B active power measurement value, [W]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(p.statusData.Em0.BActPower) },
	)
	p.collectors["c_act_power"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "c_act_power",
		Help:        "Phase C active power measurement value, [W]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(p.statusData.Em0.CActPower) },
	)

	p.collectors["total_act"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "total_act",
		Help:        "Total energy, [Wh]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(p.statusData.Emdata0.TotalAct) },
	)
	p.collectors["total_act_ret"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "shelly",
		Name:        "total_act_ret",
		Help:        "Total energy returned, [Wh]",
		ConstLabels: constLabels,
	},
		func() float64 { return float64(p.statusData.Emdata0.TotalActRet) },
	)

	var c []prometheus.Collector
	for _, v := range p.collectors {
		c = append(c, v)
	}

	return c, nil
}
