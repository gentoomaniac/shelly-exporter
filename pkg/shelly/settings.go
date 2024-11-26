package shelly

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
