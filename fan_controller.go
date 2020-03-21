package ipmi_controller

import (
	"log"
	"time"

	"github.com/ihciah/ipmi-controller/pkg/ipmi"
)

type FanController struct {
	ipmi *ipmi.IPMI
	cfg ControllerConfig
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
			switch true {
			case temp < 30:
				_, err = fc.ipmi.SetFanSpeed(8)
			case temp < 45:
				_, err = fc.ipmi.SetFanSpeed(10)
			case temp < 55:
				_, err = fc.ipmi.SetFanSpeed(13)
			case temp < 65:
				_, err = fc.ipmi.SetFanSpeed(18)
			case temp < 75:
				_, err = fc.ipmi.SetFanSpeed(25)
			case temp < 85:
				_, err = fc.ipmi.SetFanSpeed(35)
			default:
				_, err = fc.ipmi.SetFanSpeed(0)
			}
			if err != nil {
				log.Printf("error when SetFanSpeed: %v", err)
			}
		} else {
			log.Printf("error when GetTemperatureNumber: %v", err)
		}
		<-ticker.C
	}
}
