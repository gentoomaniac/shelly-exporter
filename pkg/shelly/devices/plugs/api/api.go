package api

type ActionStats struct {
	Skipped int `json:"skipped"`
}

type Cloud struct {
	Enabled   bool `json:"enabled"`
	Connected bool `json:"connected"`
}

type Meter struct {
	Counters  [3]float64 `json:"counters"`
	IsValid   bool       `json:"is_valid"`
	Power     float64    `json:"power"`
	OverPower float64    `json:"overpower"`
	Timestamp int        `json:"timestamp"`
	Total     int        `json:"total"`
}

type Mqtt struct {
	Enabled bool `json:"enabled"`
}

type Relay struct {
	HasTimer       bool   `json:"has_timer"`
	IsOn           bool   `json:"ison"`
	OverPower      bool   `json:"overpower"`
	Source         string `json:"source"`
	TimerDuration  int    `json:"time_duration"`
	TimerRemaining int    `json:"time_remaining"`
	TimerStarted   int    `json:"timer_started"`
}

type Temperature struct {
	Celsius    float32 `json:"tC"`
	Fahrenheit float32 `json:"tF"`
	IsValid    bool    `json:"is_valid"`
}

type Update struct {
	BetaVersion string `json:"beta_version"`
	HasUpdate   bool   `json:"has_update"`
	NewVersion  string `json:"new_version"`
	OldVersion  string `json:"old_version"`
	Status      string `json:"status"`
}

type Wifi struct {
	Connected bool
	IP        string
	RSSI      int16
	SSID      string
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
	Tmp               Temperature `json:"tmp"`
	Time              string      `json:"time"`
	Unixtime          int         `json:"unixtime"`
	Update            Update      `json:"update"`
	Uptime            int         `json:"uptime"`
	Wifi              Wifi        `json:"wifi_sta"`
}

type Device struct {
	Hostname   string `json:"hostname"`
	Mac        string `json:"mac"`
	NumMeters  int    `json:"num_meters"`
	NumOutputs int    `json:"num_outputs"`
	Type       string `json:"type"`
}

type Settings struct {
	Device Device `json:"device"`
	Name   string `json:"name"`
}
