package ipmi_fan_controller

import (
	"log"
	"os/exec"
)

type Executor struct {
	config
}

func NewExecutor(config config) Executor {
	return Executor{config: config}
}

func (e *Executor) GetTemperature() (int, error) {
	cmd := exec.Command(e.Executable, "-I", "lanplus", "-H", e.config.Addr, "-U", e.config.Username,
		"-P", e.config.Password, "sensor", "reading", "Temp")
	log.Println(cmd.Args)

	content, err := cmd.Output()
	if err != nil {
		log.Print(err)
		return 0, err
	}
	log.Println(string(content))
	return 0, nil
}
