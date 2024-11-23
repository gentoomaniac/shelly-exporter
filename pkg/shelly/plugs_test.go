package shelly

import (
	"encoding/json"
	"testing"
)

const (
	testResponse = "{\"wifi_sta\":{\"connected\":true,\"ssid\":\"Metalmania-iot\",\"ip\":\"10.1.3.104\",\"rssi\":-77},\"cloud\":{\"enabled\":true,\"connected\":true},\"mqtt\":{\"connected\":false},\"time\":\"17:51\",\"unixtime\":1732380718,\"serial\":14819,\"has_update\":false,\"mac\":\"80646F82988E\",\"cfg_changed_cnt\":0,\"actions_stats\":{\"skipped\":0},\"relays\":[{\"ison\":false,\"has_timer\":false,\"timer_started\":0,\"timer_duration\":0,\"timer_remaining\":0,\"overpower\":false,\"source\":\"cloud\"}],\"meters\":[{\"power\":0.00,\"overpower\":0.00,\"is_valid\":true,\"timestamp\":1732384318,\"counters\":[0.000, 0.000, 0.000],\"total\":7492}],\"temperature\":13.09,\"overtemperature\":false,\"tmp\":{\"tC\":13.09,\"tF\":55.57, \"is_valid\":true},\"update\":{\"status\":\"idle\",\"has_update\":false,\"new_version\":\"20230913-113421/v1.14.0-gcb84623\",\"old_version\":\"20230913-113421/v1.14.0-gcb84623\",\"beta_version\":\"20231107-164219/v1.14.1-rc1-g0617c15\"},\"ram_total\":52056,\"ram_free\":40436,\"fs_size\":233681,\"fs_free\":166915,\"uptime\":671403}"
)

func TestPlugSDecode(t *testing.T) {
	p := PlugS{}
	err := p.Decode([]byte(testResponse))
	if err != nil {
		t.Errorf("Decoding failed: %s", err)
	}

	blob, _ := json.Marshal(p.Status)
	t.Log(string(blob))
}
