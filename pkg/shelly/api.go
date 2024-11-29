package shelly

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
