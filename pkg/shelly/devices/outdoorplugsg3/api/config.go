package api

type Config struct {
	Ble     BleConfig    `json:"ble"`
	Bthome  any          `json:"bthome"` // Usually empty or custom
	Cloud   CloudConfig  `json:"cloud"`
	Knx     KnxConfig    `json:"knx"`
	Matter  MatterConfig `json:"matter"`
	Mqtt    MqttConfig   `json:"mqtt"`
	PlugsUI PlugsUI      `json:"plugs_ui"`
	Switch0 SwitchConfig `json:"switch:0"` // Colon handled by tag
	Sys     SysConfig    `json:"sys"`
	Wifi    WifiConfig   `json:"wifi"`
	Ws      WsConfig     `json:"ws"`
}

type BleConfig struct {
	Enable bool `json:"enable"`
	RPC    struct {
		Enable bool `json:"enable"`
	} `json:"rpc"`
}

type CloudConfig struct {
	Enable bool   `json:"enable"`
	Server string `json:"server"`
}

type KnxConfig struct {
	Enable  bool   `json:"enable"`
	IA      string `json:"ia"`
	Routing struct {
		Addr string `json:"addr"`
	} `json:"routing"`
}

type MatterConfig struct {
	Enable bool `json:"enable"`
}

type MqttConfig struct {
	Enable        bool   `json:"enable"`
	Server        string `json:"server"` // Pointer to handle null
	ClientID      string `json:"client_id"`
	User          string `json:"user"`
	TopicPrefix   string `json:"topic_prefix"`
	RPCNtf        bool   `json:"rpc_ntf"`
	StatusNtf     bool   `json:"status_ntf"`
	UseClientCert bool   `json:"use_client_cert"`
	EnableRPC     bool   `json:"enable_rpc"`
	EnableControl bool   `json:"enable_control"`
}

type PlugsUI struct {
	Leds struct {
		Mode      string         `json:"mode"`
		Colors    map[string]any `json:"colors"` // Map handles dynamic keys like switch:0
		NightMode struct {
			Enable        bool      `json:"enable"`
			Brightness    float64   `json:"brightness"`
			ActiveBetween []float64 `json:"active_between"`
		} `json:"night_mode"`
	} `json:"leds"`
	Controls map[string]any `json:"controls"`
}

type SwitchConfig struct {
	ID                       int     `json:"id"`
	Name                     *string `json:"name"`
	InitialState             string  `json:"initial_state"`
	AutoOn                   bool    `json:"auto_on"`
	AutoOnDelay              float64 `json:"auto_on_delay"`
	AutoOff                  bool    `json:"auto_off"`
	AutoOffDelay             float64 `json:"auto_off_delay"`
	PowerLimit               int     `json:"power_limit"`
	VoltageLimit             int     `json:"voltage_limit"`
	AutorecoverVoltageErrors bool    `json:"autorecover_voltage_errors"`
	CurrentLimit             float64 `json:"current_limit"`
	Reverse                  bool    `json:"reverse"`
}

type SysConfig struct {
	Device struct {
		Name         *string `json:"name"`
		Mac          string  `json:"mac"`
		FwID         string  `json:"fw_id"`
		Discoverable bool    `json:"discoverable"`
		EcoMode      bool    `json:"eco_mode"`
	} `json:"device"`
	Location struct {
		Tz  string  `json:"tz"`
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
	Sntp struct {
		Server string `json:"server"`
	} `json:"sntp"`
	CfgRev int `json:"cfg_rev"`
}

type WifiConfig struct {
	AP struct {
		SSID   string `json:"ssid"`
		Enable bool   `json:"enable"`
	} `json:"ap"`
	Sta struct {
		SSID     string  `json:"ssid"`
		Enable   bool    `json:"enable"`
		IPv4Mode string  `json:"ipv4mode"`
		IP       *string `json:"ip"`
	} `json:"sta"`
	Roam struct {
		RSSIThreshold int `json:"rssi_thr"`
		Interval      int `json:"interval"`
	} `json:"roam"`
}

type WsConfig struct {
	Enable bool    `json:"enable"`
	Server *string `json:"server"`
}
