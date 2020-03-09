package ipmi_fan_controller

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	tempNumRegexp = regexp.MustCompile(`(?m)Temp\s*\|\s*(\d*)`)
	tempRegexp    = regexp.MustCompile(`^(.*?\s*)\|.*?(\|\s+\d*\s+degrees\s+C)`)
)

type IPMI struct {
	Config
}

func NewIPMI(config Config) IPMI {
	return IPMI{Config: config}
}

func (e *IPMI) execute(arg ...string) (string, error) {
	args := []string{"-I", "lanplus", "-H", e.IPMIConfig.Addr, "-U", e.IPMIConfig.Username, "-P", e.IPMIConfig.Password}
	args = append(args, arg...)
	cmd := exec.Command(e.ControllerConfig.Executable, args...)
	log.Print("Executing", cmd.Args)
	content, err := cmd.Output()
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(content), nil
}

func (e *IPMI) GetTemperatureNumber() (int, error) {
	output, err := e.execute("sensor", "reading", "Temp")
	if err != nil {
		return 0, err
	}

	results := tempNumRegexp.FindStringSubmatch(output)
	if len(results) > 1 {
		temp, err := strconv.Atoi(results[1])
		if err != nil {
			err = fmt.Errorf("unable to parse %s", results[1])
		}
		return temp, nil
	}
	return 0, errors.New("unable to get temperature number")
}

func (e *IPMI) GetStatus() (string, error) {
	return e.execute("chassis", "status")
}

func (e *IPMI) SetPowerOn() (string, error) {
	return e.execute("power", "on")
}

func (e *IPMI) SetPowerOff() (string, error) {
	return e.execute("power", "off")
}

func (e *IPMI) GetTemperature() (string, error) {
	output, err := e.execute("sdr", "type", "Temperature")
	if err != nil {
		return "", err
	}

	lines := []string{"Temperature Sensors:"}
	for _, match := range tempRegexp.FindAllStringSubmatch(output, -1) {
		if len(match) > 2 {
			lines = append(lines, match[1]+match[2])
		}
	}
	if len(lines) <= 1 {
		return "", errors.New("unable to get temperature")
	}
	return strings.Join(lines, "\n"), nil
}

func (e *IPMI) SetFanSpeed(degree int) (string, error) {
	if degree == 0 {
		log.Print("set fan speed to auto")
		return e.execute("raw", "0x30", "0x30", "0x01", "0x01")
	}
	if degree > 100 {
		degree = 100
	}
	if degree < 0 {
		degree = 0
	}
	_, err := e.execute("raw", "0x30", "0x30", "0x01", "0x00")
	if err != nil {
		return "", err
	}

	degreeHex := fmt.Sprintf("0x%02x", degree)
	log.Printf("set fan speed to %d", degree)
	return e.execute("raw", "0x30", "0x30", "0x02", "0xff", degreeHex)
}
