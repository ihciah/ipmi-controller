package main

import (
	"flag"
	"log"

	controller "github.com/ihciah/ipmi-controller"
	"github.com/ihciah/ipmi-controller/pkg/ipmi"
)

func main() {
	var configFile = flag.String("config", "config.json", "Config file")
	flag.Parse()

	config, err := controller.NewConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	ctl := controller.NewFanController((*ipmi.IPMI)(&config.IPMIConfig), config.ControllerConfig)
	ctl.Start()
}
