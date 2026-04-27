package shelly

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"

	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/auth"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/devices/minipmg3"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/devices/outdoorplugsg3"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/devices/plugs"
	"github.com/gentoomaniac/shelly-exporter/pkg/shelly/devices/pro3em"
	"github.com/prometheus/client_golang/prometheus"
)

const TypeString = "SHELLY"

type ShellyDevice interface {
	Collectors() ([]prometheus.Collector, error)
	Hostname() string
	Name() string
	Refresh() error
	RefreshDeviceinfo() error
}

func DeviceFromIP(IP *netip.Addr, auth *auth.Auth, labels map[string]string) (ShellyDevice, error) {
	info, err := GetDeviceInfo(IP)
	if err != nil {
		return nil, err
	}

	devType := info.Type
	if devType == "" {
		devType = info.App
	}

	switch devType {
	case "OutdoorPlugSG3":
		return outdoorplugsg3.NewOutdoorPlugSG3(
			outdoorplugsg3.Config{Ip: IP, Auth: auth, Labels: labels},
		)
	case "SHPLG-S":
		return plugs.NewPlugS(
			plugs.Config{Ip: IP, Auth: auth, Labels: labels},
		)

	case "MiniPMG3":
		return minipmg3.NewMiniPMG3(
			minipmg3.Config{Ip: IP, Auth: auth, Labels: labels},
		)

	case "Pro3EM":
		return pro3em.NewPro3EM(
			pro3em.Config{Ip: IP, Auth: auth, Labels: labels},
		)
	}

	return nil, fmt.Errorf("unknown device: %s", info.ID)
}

func GetDeviceInfo(IP *netip.Addr) (*DeviceInfo, error) {
	resp, err := http.Get("http://" + IP.String() + "/shelly")
	if err != nil {
		return nil, fmt.Errorf("failed requesting device info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device info request failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading request answer: %v", err)
	}

	var devInfo DeviceInfo
	err = json.Unmarshal(body, &devInfo)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling device info: %v", err)
	}

	return &devInfo, nil
}

type DeviceInfo struct {
	// Common fields
	ID  string `json:"id,omitempty"`
	MAC string `json:"mac"`

	// Gen 1 Specific
	Auth         bool   `json:"auth,omitempty"`
	Discoverable bool   `json:"discoverable,omitempty"`
	Fw           string `json:"fw,omitempty"`
	NumOutputs   int    `json:"num_outputs,omitempty"`
	NumMeters    int    `json:"num_meters,omitempty"`
	Type         string `json:"type,omitempty"`

	// Gen 2/3 Specific
	App        string `json:"app,omitempty"`
	AuthEn     bool   `json:"auth_en,omitempty"`
	AuthDomain string `json:"auth_domain,omitempty"` // use as Realm
	FwID       string `json:"fw_id,omitempty"`
	Matter     bool   `json:"matter,omitempty"`
	Name       string `json:"name,omitempty"`
	Profile    string `json:"profile,omitempty"`
	Slot       int    `json:"slot,omitempty"`
	Gen        int    `json:"gen,omitempty"` // 0 or missing for Gen 1
	Model      string `json:"model,omitempty"`
	Ver        string `json:"ver,omitempty"`
}
