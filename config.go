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

type SpeedConfig struct {
	Temperature int `json:"temperature"`
	Percentage  int `json:"percentage"`
}

type ControllerConfig struct {
	Ticker      int           `json:"ticker"`
	SpeedConfig []SpeedConfig `json:"speed"`
}

func (c *ControllerConfig) Validate() error {
	if c.Ticker == 0 {
		c.Ticker = 60
	}
	if c.SpeedConfig == nil || len(c.SpeedConfig) == 0 {
		c.SpeedConfig = []SpeedConfig{
			{40, 1},
			{45, 2},
			{55, 3},
			{65, 5},
			{75, 10},
			{85, 20},
			{90, 40},
		}
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
