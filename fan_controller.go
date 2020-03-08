package ipmi_fan_controller

import (
	"context"
	"errors"
	"log"

	"github.com/ihciah/bmc"
)

type fanController struct {
	ctx  context.Context
	sess bmc.Session
}

func NewFanController() *fanController {
	return &fanController{}
}

func (fc *fanController) Execute(ctx context.Context, sess bmc.Session) error {
	fc.ctx, fc.sess = ctx, sess
	temperature, err := fc.getTemperature()
	if err != nil {
		log.Print(err)
		return nil
	}
	log.Printf("temperature: %v", temperature)
	return nil
}

func (fc *fanController) getTemperature() (float64, error) {
	repo, err := bmc.RetrieveSDRRepository(fc.ctx, fc.sess)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	var temperature float64
	for _, v := range repo {
		// Will filter out maximum temperature for cpus
		if v.Identity != "Temp" {
			continue
		}
		reader, err := bmc.NewSensorReader(v)
		if err != nil {
			continue
		}
		reading, err := reader.Read(fc.ctx, fc.sess)
		if reading > temperature {
			temperature = reading
		}
		log.Printf("temperature from sonsor: %v, %v%v", v.Identity, reading, v.BaseUnit.Symbol())
	}
	if temperature == 0 {
		return temperature, errors.New("unable to find sensor")
	}
	return temperature, nil
}
