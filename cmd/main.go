package main

import (
	"log"

	controller "github.com/ihciah/ipmi-fan-controller"
)

func main() {
	config, err := controller.NewConfigFromFile("config.json")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
	ipmi := controller.NewIPMI(config)

	controller.ControlFanSpeed(&ipmi)
}
