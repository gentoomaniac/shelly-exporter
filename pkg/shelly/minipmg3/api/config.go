package api

type Config struct {
	Ble struct {
		Enable bool `json:"enable"`
		RPC    struct {
			Enable bool `json:"enable"`
		} `json:"rpc"`
		Observer struct {
			Enable bool `json:"enable"`
		} `json:"observer"`
	} `json:"ble"`
	Bthome struct{} `json:"bthome"`
	Cloud  struct {
		Enable bool   `json:"enable"`
		Server string `json:"server"`
	} `json:"cloud"`
	Mqtt struct {
		Enable        bool        `json:"enable"`
		Server        interface{} `json:"server"`
		ClientID      string      `json:"client_id"`
		User          interface{} `json:"user"`
		SslCa         interface{} `json:"ssl_ca"`
		TopicPrefix   string      `json:"topic_prefix"`
		RPCNtf        bool        `json:"rpc_ntf"`
		StatusNtf     bool        `json:"status_ntf"`
		UseClientCert bool        `json:"use_client_cert"`
		EnableRPC     bool        `json:"enable_rpc"`
		EnableControl bool        `json:"enable_control"`
	} `json:"mqtt"`
	Pm10 struct {
		ID      int         `json:"id"`
		Name    interface{} `json:"name"`
		Reverse bool        `json:"reverse"`
	} `json:"pm1:0"`
	Sys struct {
		Device struct {
			Name         string `json:"name"`
			Mac          string `json:"mac"`
			FwID         string `json:"fw_id"`
			Discoverable bool   `json:"discoverable"`
			EcoMode      bool   `json:"eco_mode"`
		} `json:"device"`
		Location struct {
			Tz  string  `json:"tz"`
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"location"`
		Debug struct {
			Level     int         `json:"level"`
			FileLevel interface{} `json:"file_level"`
			Mqtt      struct {
				Enable bool `json:"enable"`
			} `json:"mqtt"`
			Websocket struct {
				Enable bool `json:"enable"`
			} `json:"websocket"`
			UDP struct {
				Addr interface{} `json:"addr"`
			} `json:"udp"`
		} `json:"debug"`
		UIData struct{} `json:"ui_data"`
		RPCUDP struct {
			DstAddr    interface{} `json:"dst_addr"`
			ListenPort interface{} `json:"listen_port"`
		} `json:"rpc_udp"`
		Sntp struct {
			Server string `json:"server"`
		} `json:"sntp"`
		CfgRev int `json:"cfg_rev"`
	} `json:"sys"`
	Wifi struct {
		Ap struct {
			Ssid          string `json:"ssid"`
			IsOpen        bool   `json:"is_open"`
			Enable        bool   `json:"enable"`
			RangeExtender struct {
				Enable bool `json:"enable"`
			} `json:"range_extender"`
		} `json:"ap"`
		Sta struct {
			Ssid       string      `json:"ssid"`
			IsOpen     bool        `json:"is_open"`
			Enable     bool        `json:"enable"`
			Ipv4Mode   string      `json:"ipv4mode"`
			IP         interface{} `json:"ip"`
			Netmask    interface{} `json:"netmask"`
			Gw         interface{} `json:"gw"`
			Nameserver interface{} `json:"nameserver"`
		} `json:"sta"`
		Sta1 struct {
			Ssid       interface{} `json:"ssid"`
			IsOpen     bool        `json:"is_open"`
			Enable     bool        `json:"enable"`
			Ipv4Mode   string      `json:"ipv4mode"`
			IP         interface{} `json:"ip"`
			Netmask    interface{} `json:"netmask"`
			Gw         interface{} `json:"gw"`
			Nameserver interface{} `json:"nameserver"`
		} `json:"sta1"`
		Roam struct {
			RssiThr  int `json:"rssi_thr"`
			Interval int `json:"interval"`
		} `json:"roam"`
	} `json:"wifi"`
	Ws struct {
		Enable bool        `json:"enable"`
		Server interface{} `json:"server"`
		SslCa  string      `json:"ssl_ca"`
	} `json:"ws"`
}
