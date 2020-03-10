package main

import (
	"flag"
	"log"

	controller "github.com/ihciah/ipmi-controller"
)

func main() {
	var configFile = flag.String("config", "config.json", "Config file")
	flag.Parse()

	config, err := controller.NewConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
	ipmi := controller.NewIPMI(config)
	go func() {
		telegramBot, err := controller.NewTelegramBot(ipmi)
		if err != nil {
			log.Printf("can not run telegram bot: %v", err)
		} else {
			telegramBot.Serve()
		}
	}()
	controller.ControlFanSpeed(ipmi)
}
