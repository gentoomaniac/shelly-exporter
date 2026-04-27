package components

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/auth"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/request"
)

type BTHomeComponent struct {
	Key     string
	Status  *BTHomeDeviceStatus
	Config  *BTHomeDeviceConfig
	Attrs   *BTHomeDeviceAttrs
	Sensors []BTHomeSensor
}

type BTHomeSensor struct {
	Status *BTHomeSensorStatus
	Config *BTHomeSensorConfig
}

type BTHomeComponents map[string]*BTHomeComponent

func GetBTHomeComponents(baseUrl *url.URL, auth *auth.Auth) (BTHomeComponents, error) {
	componentsUrl := baseUrl.JoinPath("Shelly.GetComponents")
	resp, err := request.DigestAuthedRequest(componentsUrl, auth, map[string]string{"id": "0"})
	if err != nil {
		return nil, err
	}

	return parseBTHomeComponents(resp)

}

func parseBTHomeComponents(data []byte) (BTHomeComponents, error) {
	btHomeComponents := make(map[string]*BTHomeComponent)

	var response Response
	err := json.Unmarshal(data, &response)
	if err != nil {
		return btHomeComponents, err
	}

	ensureEntry := func(addr string) *BTHomeComponent {
		if _, ok := btHomeComponents[addr]; !ok {
			btHomeComponents[addr] = &BTHomeComponent{
				Sensors: []BTHomeSensor{}, // Pre-allocate slice
			}
		}
		return btHomeComponents[addr]
	}

	for _, c := range response.Components {
		if strings.HasPrefix(c.Key, "bthomedevice:") {
			var config BTHomeDeviceConfig
			if err := json.Unmarshal(c.Config, &config); err != nil {
				return btHomeComponents, err
			}

			var status BTHomeDeviceStatus
			if err := json.Unmarshal(c.Status, &status); err != nil {
				return btHomeComponents, err
			}

			var attrs BTHomeDeviceAttrs
			if err := json.Unmarshal(c.Attrs, &attrs); err != nil {
				return btHomeComponents, err
			}

			comp := ensureEntry(config.Addr)
			comp.Key = c.Key
			comp.Config = &config
			comp.Status = &status
			comp.Attrs = &attrs

		} else if strings.HasPrefix(c.Key, "bthomesensor:") {
			var status BTHomeSensorStatus
			if err := json.Unmarshal(c.Status, &status); err != nil {
				return btHomeComponents, err
			}

			var config BTHomeSensorConfig
			if err := json.Unmarshal(c.Config, &config); err != nil {
				return btHomeComponents, err
			}

			sensor := BTHomeSensor{Config: &config, Status: &status}
			comp := ensureEntry(config.Addr)
			comp.Sensors = append(comp.Sensors, sensor)
		}
	}

	return btHomeComponents, nil
}
