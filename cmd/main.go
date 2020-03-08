package main

import (
	controller "github.com/ihciah/ipmi-fan-controller"
)

func main() {
	config, err := controller.NewConfigFromFile("config.json")
	if err != nil {
		panic(err)
	}
	controller.ExecuteLoop(config, controller.NewFanController())
}
