package v1

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

const TypeString = "HWE-P1"

func NewP1(c Config) (*Homewizard, error) {
	if c.BaseUrl == nil {
		if c.IP == nil {
			return nil, fmt.Errorf("must provide at least one of `Baseurl` or `IP`")
		}
		c.BaseUrl, _ = url.Parse(fmt.Sprintf("http://%s/", c.IP.String()))
	}
	h := &Homewizard{
		config:     c,
		collectors: make(map[string]prometheus.Collector),
	}

	if err := h.RefreshDeviceinfo(); err != nil {
		return nil, err
	}

	return h, nil
}

type Config struct {
	BaseUrl *url.URL
	Labels  map[string]string
	IP      net.IP
}

type Homewizard struct {
	config Config

	info Info
	data Data

	collectors map[string]prometheus.Collector
}

func (h Homewizard) Name() string {
	return h.info.ProductName
}

func (h *Homewizard) RefreshDeviceinfo() error {
	infoUrl := h.config.BaseUrl.JoinPath(infoEndpoint)
	resp, err := request(infoUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &h.info)
	if err != nil {
		return err
	}

	return nil
}

func (h *Homewizard) Refresh() error {
	dataUrl := h.config.BaseUrl.JoinPath(dataEndpoint)
	resp, err := request(dataUrl)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &h.data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Homewizard) Collectors() ([]prometheus.Collector, error) {
	constLabels := prometheus.Labels{
		"type":   TypeString,
		"serial": h.info.Serial,
		"name":   h.info.ProductName,
	}

	for k, v := range h.config.Labels {
		constLabels[k] = v
	}

	h.collectors["active_power_w"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "homewizard",
		Name:        "active_power_w",
		Help:        "Current real AC power being drawn, in Watts",
		ConstLabels: constLabels,
	},
		func() float64 { return h.data.ActivePowerW },
	)
	h.collectors["total_power_import_kwh"] = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:   "homewizard",
		Name:        "total_power_import_kwh",
		Help:        "Total imported power in kWh",
		ConstLabels: constLabels,
	},
		func() float64 { return h.data.ActivePowerW },
	)

	var c []prometheus.Collector
	for _, v := range h.collectors {
		c = append(c, v)
	}

	return c, nil
}
