package main

import (
	"log"

	controller "github.com/ihciah/ipmi-fan-controller"
)

func main() {
	config, err := controller.NewConfigFromFile("config.json")
	if err != nil {
		panic(err)
	}
	executor := controller.NewExecutor(config)
	_, err = executor.GetTemperature()
	if err != nil {
		log.Print(err)
	}
}
