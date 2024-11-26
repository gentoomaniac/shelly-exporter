package exporter

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"
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

func getEnv(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return value
}

type EnvString string

func (e *EnvString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	err := unmarshal(&str)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`\${env:(?P<Env>[a-zA-Z0-9-_]+)(:-(?P<Default>.+))?}`)
	matches := r.FindStringSubmatch(str)
	if len(matches) > 0 {
		envName := matches[1]
		defaultValue := ""
		if len(matches) == 4 {
			defaultValue = matches[3]
		}
		log.Debug().Str("envName", envName).Str("default", defaultValue).Msg("")

		*e = EnvString(getEnv(envName, defaultValue))
	} else {
		*e = EnvString(str)
	}

	//	return errors.New(fmt.Sprintf("invalid type: %s", typeString))

	return nil
}

type Device struct {
	Type      Type              `yaml:"type"`
	Ip        net.IP            `yaml:"ip"`
	User      EnvString         `yaml:"user"`
	Password  EnvString         `yaml:"password"`
	Frequency EnvString         `yaml:"frequency"`
	Labels    map[string]string `yaml:"labels"`
}

type Global struct {
	User      EnvString `yaml:"user"`
	Password  EnvString `yaml:"password"`
	Frequency EnvString `yaml:"frequency"`
}

type Config struct {
	Global  Global   `yaml:"global"`
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
