package api

type Status struct {
	Ble    struct{} `json:"ble"`
	Bthome struct {
		Errors []string `json:"errors"`
	} `json:"bthome"`
	Cloud struct {
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Knx    struct{} `json:"knx"`
	Matter struct {
		NumFabrics     int  `json:"num_fabrics"`
		Commissionable bool `json:"commissionable"`
	} `json:"matter"`
	Mqtt struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
	PlugsUI struct{} `json:"plugs_ui"`
	// Switch:0 contains the core power metrics
	Switch0 struct {
		ID      int     `json:"id"`
		Source  string  `json:"source"`
		Output  bool    `json:"output"`
		APower  float64 `json:"apower"`
		Voltage float64 `json:"voltage"`
		Freq    float64 `json:"freq"`
		Current float64 `json:"current"`
		AEnergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int64     `json:"minute_ts"`
		} `json:"aenergy"`
		RetAEnergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int64     `json:"minute_ts"`
		} `json:"ret_aenergy"`
		Temperature struct {
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		} `json:"temperature"`
	} `json:"switch:0"`
	Sys struct {
		Mac              string `json:"mac"`
		RestartRequired  bool   `json:"restart_required"`
		Time             string `json:"time"`
		UnixTime         int64  `json:"unixtime"`
		LastSyncTs       int64  `json:"last_sync_ts"`
		Uptime           int64  `json:"uptime"`
		RamSize          int    `json:"ram_size"`
		RamFree          int    `json:"ram_free"`
		RamMinFree       int    `json:"ram_min_free"`
		FsSize           int    `json:"fs_size"`
		FsFree           int    `json:"fs_free"`
		CfgRev           int    `json:"cfg_rev"`
		KvsRev           int    `json:"kvs_rev"`
		ScheduleRev      int    `json:"schedule_rev"`
		WebhookRev       int    `json:"webhook_rev"`
		BtrelayRev       int    `json:"btrelay_rev"`
		BthcRev          int    `json:"bthc_rev"`
		AvailableUpdates struct {
			Beta struct {
				Version string `json:"version"`
			} `json:"beta"`
			Stable struct {
				Version string `json:"version"`
			} `json:"stable"`
		} `json:"available_updates"`
		ResetReason int `json:"reset_reason"`
		UtcOffset   int `json:"utc_offset"`
	} `json:"sys"`
	Wifi struct {
		StaIP  string   `json:"sta_ip"`
		Status string   `json:"status"`
		SSID   string   `json:"ssid"`
		BSSID  string   `json:"bssid"`
		RSSI   int      `json:"rssi"`
		StaIP6 []string `json:"sta_ip6"`
	} `json:"wifi"`
	Ws struct {
		Connected bool `json:"connected"`
	} `json:"ws"`
}
