package exporter

import (
	"errors"
	"fmt"
	"net"

	"gopkg.in/yaml.v2"
)

const (
	SHPLG_S Type = iota
)

type Type uint8

var (
	TypeString = map[Type]string{
		SHPLG_S: "SHPLG-S",
	}
	StringType = map[string]Type{
		"SHPLG-S": SHPLG_S,
	}
)

func (t *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var typeString string
	err := unmarshal(&typeString)
	if err != nil {
		return err
	}

	var ok bool
	*t, ok = StringType[typeString]
	if !ok {
		return errors.New(fmt.Sprintf("invalid type: %s", typeString))
	}

	return nil
}

func (t *Type) MarshalYAML() (interface{}, error) {
	return TypeString[*t], nil
}

type Device struct {
	Name      string `yaml:"name"`
	Alias     string `yaml:"alias"`
	Type      Type   `yaml:"type"`
	Ip        net.IP `yaml:"ip"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Frequency string `yaml:"frequency"`
}

type Config struct {
	Devices []Device `yaml:"devices"`
}

func NewConfigFromContent(filecontent []byte) (*Config, error) {
	c := &Config{}

	err := yaml.Unmarshal(filecontent, c)
	if err != nil {
		return nil, fmt.Errorf("Unmarshaling config failed: %w", err)
	}

	return c, nil
}
