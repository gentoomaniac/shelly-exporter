package api

type Status struct {
	Ble    struct{} `json:"ble"`
	Bthome struct {
		Errors []string `json:"errors"`
	} `json:"bthome"`
	Cloud struct {
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Em0 struct {
		ID                  int           `json:"id"`
		ACurrent            float64       `json:"a_current"`
		AVoltage            float64       `json:"a_voltage"`
		AActPower           float64       `json:"a_act_power"`
		AAprtPower          float64       `json:"a_aprt_power"`
		APf                 float64       `json:"a_pf"`
		AFreq               float64       `json:"a_freq"`
		BCurrent            float64       `json:"b_current"`
		BVoltage            float64       `json:"b_voltage"`
		BActPower           float64       `json:"b_act_power"`
		BAprtPower          float64       `json:"b_aprt_power"`
		BPf                 float64       `json:"b_pf"`
		BFreq               float64       `json:"b_freq"`
		CCurrent            float64       `json:"c_current"`
		CVoltage            float64       `json:"c_voltage"`
		CActPower           float64       `json:"c_act_power"`
		CAprtPower          float64       `json:"c_aprt_power"`
		CPf                 float64       `json:"c_pf"`
		CFreq               float64       `json:"c_freq"`
		NCurrent            interface{}   `json:"n_current"`
		TotalCurrent        float64       `json:"total_current"`
		TotalActPower       float64       `json:"total_act_power"`
		TotalAprtPower      float64       `json:"total_aprt_power"`
		UserCalibratedPhase []interface{} `json:"user_calibrated_phase"`
	} `json:"em:0"`
	Emdata0 struct {
		ID                 int     `json:"id"`
		ATotalActEnergy    float64 `json:"a_total_act_energy"`
		ATotalActRetEnergy float64 `json:"a_total_act_ret_energy"`
		BTotalActEnergy    float64 `json:"b_total_act_energy"`
		BTotalActRetEnergy float64 `json:"b_total_act_ret_energy"`
		CTotalActEnergy    float64 `json:"c_total_act_energy"`
		CTotalActRetEnergy float64 `json:"c_total_act_ret_energy"`
		TotalAct           float64 `json:"total_act"`
		TotalActRet        float64 `json:"total_act_ret"`
	} `json:"emdata:0"`
	Eth struct {
		IP interface{} `json:"ip"`
	} `json:"eth"`
	Modbus struct{} `json:"modbus"`
	Mqtt   struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
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
	Temperature0 struct {
		ID int     `json:"id"`
		TC float64 `json:"tC"`
		TF float64 `json:"tF"`
	} `json:"temperature:0"`
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
