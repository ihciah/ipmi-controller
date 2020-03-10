package ipmi_controller

import (
	"log"
	"time"
)

func ControlFanSpeed(e *IPMI) {
	ticker := time.NewTicker(time.Duration(e.Config.ControllerConfig.Ticker) * time.Second)
	defer ticker.Stop()
	for {
		temp, err := e.GetTemperatureNumber()
		if err == nil {
			log.Printf("get temperature %d", temp)
			switch true {
			case temp < 30:
				_, err = e.SetFanSpeed(10)
			case temp < 45:
				_, err = e.SetFanSpeed(13)
			case temp < 55:
				_, err = e.SetFanSpeed(16)
			case temp < 65:
				_, err = e.SetFanSpeed(20)
			case temp < 75:
				_, err = e.SetFanSpeed(30)
			case temp < 85:
				_, err = e.SetFanSpeed(50)
			default:
				_, err = e.SetFanSpeed(0)
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
