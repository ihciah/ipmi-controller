package ipmi_controller

import (
	"log"
	"time"

	"github.com/ihciah/ipmi-controller/pkg/ipmi"
)

type FanController struct {
	ipmi *ipmi.IPMI
	cfg  ControllerConfig
}

func NewFanController(ipmi *ipmi.IPMI, cfg ControllerConfig) *FanController {
	return &FanController{
		ipmi: ipmi,
		cfg:  cfg,
	}
}

func (fc *FanController) Start() {
	ticker := time.NewTicker(time.Duration(fc.cfg.Ticker) * time.Second)
	defer ticker.Stop()
	for {
		temp, err := fc.ipmi.GetTemperatureNumber()
		if err == nil {
			log.Printf("get temperature %d", temp)
			_, err = fc.ipmi.SetFanSpeed(getSpeed(temp, fc.cfg.SpeedConfig))
			if err != nil {
				log.Printf("error when SetFanSpeed: %v", err)
			}
		} else {
			log.Printf("error when GetTemperatureNumber: %v", err)
		}
		<-ticker.C
	}
}

func getSpeed(temp int, config []SpeedConfig) int {
	for _, pair := range config {
		if temp < pair.Temperature {
			return pair.Percentage
		}
	}
	return 0
}
