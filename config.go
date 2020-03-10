package ipmi_controller

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	IPMIConfig       IPMIConfig       `json:"ipmi"`
	ControllerConfig ControllerConfig `json:"controller"`
	TelegramConfig   TelegramConfig   `json:"telegram"`
}
type IPMIConfig struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type ControllerConfig struct {
	Ticker     int    `json:"ticker"`
	Executable string `json:"executable"`
}
type TelegramConfig struct {
	Token       string `json:"token"`
	URL         string `json:"url"`
	Admin       []int  `json:"admin"`
	PollTimeout int    `json:"poll_timeout"`
}

func (c *IPMIConfig) validate() error {
	if c.Addr == "" {
		return errors.New("IPMI addr can not be blank")
	}
	if c.Username == "" {
		return errors.New("IPMI username can not be blank")
	}
	if c.Password == "" {
		return errors.New("IPMI password can not be blank")
	}
	return nil
}

func (c *ControllerConfig) validate() error {
	if c.Ticker == 0 {
		c.Ticker = 60
	}
	if c.Executable == "" {
		c.Executable = "ipmitool"
	}
	return nil
}

func (c *TelegramConfig) validate() error {
	return nil
}

func (c *Config) validate() (err error) {
	if err = c.IPMIConfig.validate(); err != nil {
		return
	}
	if err = c.ControllerConfig.validate(); err != nil {
		return
	}
	if err = c.TelegramConfig.validate(); err != nil {
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
	return c, c.validate()
}
