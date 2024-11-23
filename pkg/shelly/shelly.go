package shelly

type Meter struct {
	Power     float64
	OverPower float64
	IsValid   bool
	Timestamp int
	Counters  [3]float64
}

type Relay struct {
	IsOn           bool
	HasTimer       bool
	TimerStarted   int
	TimerDuration  int
	TimerRemaining int
	OverPower      bool
}

type Temperature struct {
	Celsius         float32
	Fahrnheit       float32
	IsValid         bool
	OverTemperature bool
}

type Wifi struct {
	Connected bool
	SSID      string
	IP        string
	RSSI      int16
}

type Update struct {
	Status      string
	HasUpdate   bool
	NewVersion  string
	OldVersion  string
	BetaVersion string
}

type System struct {
	Timestamp int
	Serial    int
	MAC       string

	RamTotal int
	RamFree  int
	FSSize   int
	FSFree   int
	Uptime   int
}
