package api

type Status struct {
	Ble    struct{} `json:"ble"`
	Bthome struct {
		Errors []string `json:"errors"`
	} `json:"bthome"`
	Cloud struct {
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Mqtt struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
	Pm10 struct {
		ID      int     `json:"id"`
		Voltage float64 `json:"voltage"`
		Current float64 `json:"current"`
		Apower  float64 `json:"apower"`
		Freq    float64 `json:"freq"`
		Aenergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int       `json:"minute_ts"`
		} `json:"aenergy"`
		RetAenergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int       `json:"minute_ts"`
		} `json:"ret_aenergy"`
	} `json:"pm1:0"`
	Sys struct {
		Mac              string   `json:"mac"`
		RestartRequired  bool     `json:"restart_required"`
		Time             string   `json:"time"`
		Unixtime         int      `json:"unixtime"`
		Uptime           int      `json:"uptime"`
		RAMSize          int      `json:"ram_size"`
		RAMFree          int      `json:"ram_free"`
		FsSize           int      `json:"fs_size"`
		FsFree           int      `json:"fs_free"`
		CfgRev           int      `json:"cfg_rev"`
		KvsRev           int      `json:"kvs_rev"`
		ScheduleRev      int      `json:"schedule_rev"`
		WebhookRev       int      `json:"webhook_rev"`
		AvailableUpdates struct{} `json:"available_updates"`
		ResetReason      int      `json:"reset_reason"`
	} `json:"sys"`
	Wifi struct {
		StaIP  string `json:"sta_ip"`
		Status string `json:"status"`
		Ssid   string `json:"ssid"`
		Rssi   int    `json:"rssi"`
	} `json:"wifi"`
	Ws struct {
		Connected bool `json:"connected"`
	} `json:"ws"`
}
