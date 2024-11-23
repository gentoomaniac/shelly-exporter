package shelly

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
)

func NewPlugS(ip net.IP, user string, password string) PlugS {
	baseUrl, _ := url.Parse(fmt.Sprintf("http://%s/", ip.String()))
	return PlugS{
		baseUrl: baseUrl,
		auth:    Auth{user: user, password: password},
	}
}

type PlugS struct {
	baseUrl *url.URL
	auth    Auth

	Status Status
}

type Status struct {
	ActionStats       ActionStats `json:"action_stats"`
	ConfigChangeCount int         `json:"cfg_change_count"`
	Cloud             Cloud       `json:"cloud"`
	FsFree            int         `json:"fs_free"`
	FsSize            int         `json:"fs_size"`
	HasUpdate         bool        `json:"has_update"`
	Mac               string      `json:"mac"`
	Meters            []Meter     `json:"meters"`
	Mqtt              Mqtt        `json:"mqtt"`
	OverTemperature   bool        `json:"overtemperature"`
	RamFree           int         `json:"ram_free"`
	RamTotal          int         `json:"ram_total"`
	Relays            []Relay     `json:"relays"`
	Serial            int         `json:"serial"`
	Temperature       float32     `json:"temperature"`
	Time              string      `json:"time"`
	Unixtime          int         `json:"unixtime"`
	Update            Update      `json:"update"`
	Uptime            int         `json:"uptime"`
	Wifi              Wifi        `json:"wifi_sta"`
}

func (p *PlugS) Decode(rawResponse []byte) error {
	err := json.Unmarshal(rawResponse, &p.Status)
	if err != nil {
		return err
	}

	return nil
}

func (p *PlugS) Update() error {
	statusUrl := p.baseUrl.JoinPath("status")
	resp, err := request(statusUrl, &p.auth)
	if err != nil {
		return err
	}

	return p.Decode(resp)
}
