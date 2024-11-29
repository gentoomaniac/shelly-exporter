package v1

var infoEndpoint = "/api"

type Info struct {
	ProductType     string `json:"product_type"`
	ProductName     string `json:"product_name"`
	Serial          string `json:"serial"`
	FirmwareVersion string `json:"firmware_version"`
	APIVersion      string `json:"api_version"`
}

var dataEndpoint = "/api/v1/data"

type Data struct {
	WifiSsid                 string  `json:"wifi_ssid"`
	WifiStrength             int     `json:"wifi_strength"`
	SmrVersion               int     `json:"smr_version"`
	MeterModel               string  `json:"meter_model"`
	UniqueID                 string  `json:"unique_id"`
	ActiveTariff             int     `json:"active_tariff"`
	TotalPowerImportKwh      float64 `json:"total_power_import_kwh"`
	TotalPowerImportT1Kwh    float64 `json:"total_power_import_t1_kwh"`
	TotalPowerImportT2Kwh    float64 `json:"total_power_import_t2_kwh"`
	TotalPowerExportKwh      float64 `json:"total_power_export_kwh"`
	TotalPowerExportT1Kwh    float64 `json:"total_power_export_t1_kwh"`
	TotalPowerExportT2Kwh    float64 `json:"total_power_export_t2_kwh"`
	ActivePowerW             float64 `json:"active_power_w"`
	ActivePowerL1W           float64 `json:"active_power_l1_w"`
	ActivePowerL2W           float64 `json:"active_power_l2_w"`
	ActivePowerL3W           float64 `json:"active_power_l3_w"`
	ActiveCurrentL1A         float64 `json:"active_current_l1_a"`
	ActiveCurrentL2A         float64 `json:"active_current_l2_a"`
	ActiveCurrentL3A         float64 `json:"active_current_l3_a"`
	VoltageSagL1Count        float64 `json:"voltage_sag_l1_count"`
	VoltageSagL2Count        float64 `json:"voltage_sag_l2_count"`
	VoltageSagL3Count        float64 `json:"voltage_sag_l3_count"`
	VoltageSwellL1Count      float64 `json:"voltage_swell_l1_count"`
	VoltageSwellL2Count      float64 `json:"voltage_swell_l2_count"`
	VoltageSwellL3Count      float64 `json:"voltage_swell_l3_count"`
	AnyPowerFailCount        int     `json:"any_power_fail_count"`
	LongPowerFailCount       int     `json:"long_power_fail_count"`
	TotalGasM3               float64 `json:"total_gas_m3"`
	GasTimestamp             int64   `json:"gas_timestamp"`
	GasUniqueID              string  `json:"gas_unique_id"`
	ActivePowerAverageW      float64 `json:"active_power_average_w"`
	MontlyPowerPeakW         float64 `json:"montly_power_peak_w"`
	MontlyPowerPeakTimestamp int64   `json:"montly_power_peak_timestamp"`
	External                 []struct {
		UniqueID  string  `json:"unique_id"`
		Type      string  `json:"type"`
		Timestamp int64   `json:"timestamp"`
		Value     float64 `json:"value"`
		Unit      string  `json:"unit"`
	} `json:"external"`
}
