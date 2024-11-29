package config

import (
	"testing"
)

const (
	TestConfigData = `
devices:
- name: "shellyplug-s-80646F819FD8"
  alias: "birds"
  type: "SHPLG-S"
  ip: "10.1.3.117"
  user: "${env:SHELLY_USER:-marco}"
  password: "${env:SHELLY_PASSWORD}"
`
)

func TestConfig(t *testing.T) {
	_, err := NewConfigFromContent([]byte(TestConfigData))
	if err != nil {
		t.Errorf("Decoding failed: %s", err)
	}
}
