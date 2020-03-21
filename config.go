package ipmi_controller

import (
	"encoding/json"
	"os"

	"github.com/ihciah/ipmi-controller/pkg/ipmi"
)

type Config struct {
	IPMIConfig       ipmi.IPMIConfig  `json:"ipmi"`
	ControllerConfig ControllerConfig `json:"controller"`
}

type ControllerConfig struct {
	Ticker int `json:"ticker"`
}

func (c *ControllerConfig) Validate() error {
	if c.Ticker == 0 {
		c.Ticker = 60
	}
	return nil
}

func (c *Config) Validate() (err error) {
	if err = c.IPMIConfig.Validate(); err != nil {
		return
	}
	if err = c.ControllerConfig.Validate(); err != nil {
		return
	}
	return
}

func NewConfigFromFile(path string) (Config, error) {
	var c Config
	file, err := os.Open(path)
	if err != nil {
		return c, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return c, err
	}
	return c, c.Validate()
}
