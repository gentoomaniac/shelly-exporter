package components

import "encoding/json"

type ObjectID int

const (
	Battery     ObjectID = 1
	Humidity    ObjectID = 46
	Temperature ObjectID = 69
)

func (o ObjectID) String() string {
	switch o {
	case Battery:
		return "battery"
	case Humidity:
		return "humidity"
	case Temperature:
		return "temperature"
	default:
		return "unknown"
	}
}

type Response struct {
	Components []Component `json:"components"`
}

type Component struct {
	Key    string          `json:"key"`
	Status json.RawMessage `json:"status"`
	Config json.RawMessage `json:"config"`
	Attrs  json.RawMessage `json:"attrs,omitempty"`
}

type BTHomeDeviceStatus struct {
	ID            int    `json:"id"`
	RSSI          int    `json:"rssi"`
	Battery       int    `json:"battery"`
	PacketID      int    `json:"packet_id"`
	LastUpdatedTs int64  `json:"last_updated_ts"`
	Paired        bool   `json:"paired"`
	RPC           bool   `json:"rpc"`
	RSV           int    `json:"rsv"`
	FwVer         string `json:"fw_ver"`
}

type BTHomeDeviceConfig struct {
	ID   int     `json:"id"`
	Addr string  `json:"addr"`
	Name string  `json:"name"`
	Key  *string `json:"key"`
	Meta struct {
		UI struct {
			View      string  `json:"view"`
			LocalName string  `json:"local_name"`
			Icon      *string `json:"icon"`
		} `json:"ui"`
	} `json:"meta"`
}

type BTHomeDeviceAttrs struct {
	Flags   int `json:"flags"`
	ModelID int `json:"model_id"`
}

type BTHomeSensorStatus struct {
	ID            int     `json:"id"`
	Value         float64 `json:"value"`
	LastUpdatedTs int64   `json:"last_updated_ts"`
}

type BTHomeSensorConfig struct {
	ID    int      `json:"id"`
	Addr  string   `json:"addr"`
	Name  *string  `json:"name"`
	ObjID ObjectID `json:"obj_id"`
	Idx   int      `json:"idx"`
	Meta  *struct {
		UI struct {
			Icon *string `json:"icon"`
		} `json:"ui"`
	} `json:"meta"`
}
