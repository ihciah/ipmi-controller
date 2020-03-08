package ipmi_fan_controller

import (
	"encoding/json"
	"errors"
	"os"
)

type config struct {
	Addr       string `json:"addr"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Timeout    int    `json:"timeout"`
	Ticker     int    `json:"ticker"`
	Executable string `json:"executable"`
}

func (c *config) validate() error {
	if c.Timeout <= 0 {
		return errors.New("timeout cannot be 0")
	}
	if c.Ticker <= 5 {
		return errors.New("ticker cannot be less than 5")
	}
	return nil
}

func NewConfig(executable, addr, username, password string, timeout, ticker int) (config, error) {
	c := config{
		Addr:       addr,
		Username:   username,
		Password:   password,
		Timeout:    timeout,
		Ticker:     ticker,
		Executable: executable,
	}
	return c, c.validate()
}

func NewConfigFromFile(path string) (config, error) {
	var c config
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
