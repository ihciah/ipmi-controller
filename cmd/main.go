package main

import (
	"flag"
	"log"
	"time"

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
		for {
			telegramBot, err := controller.NewTelegramBot(ipmi)
			if err != nil {
				log.Printf("can not run telegram bot: %v", err)
				if ipmi.TelegramConfig.Token == "" {
					log.Print("telegram is not enabled")
					return
				}
				log.Print("will sleep 10 seconds and retry starting telegram bot")
				time.Sleep(10 * time.Second)
			} else {
				telegramBot.Serve()
			}
		}
	}()
	controller.ControlFanSpeed(ipmi)
}
