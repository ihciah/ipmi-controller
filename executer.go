package ipmi_fan_controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ihciah/bmc"
	"github.com/ihciah/bmc/pkg/ipmi"
)

type config struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
	Ticker   int    `json:"ticker"`
}

func (c *config) validate() error {
	if c.Timeout <= 0 {
		return errors.New("timeout cannot be 0")
	}
	if c.Ticker <= 5 {
		return errors.New("ticker cannot be less than 5")
	}
	return nil
}

func NewConfig(addr, username, password string, timeout, ticker int) (config, error) {
	c := config{
		Addr:     addr,
		Username: username,
		Password: password,
		Timeout:  timeout,
		Ticker:   ticker,
	}
	return c, c.validate()
}

func NewConfigFromFile(path string) (config, error) {
	var c config
	file, err := os.Open(path)
	if err != nil {
		return c, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return c, err
	}
	return c, c.validate()
}

type Task interface {
	Execute(ctx context.Context, sess bmc.Session) error
}

func Execute(config config, task Task) error {
	machine, err := bmc.DialV2(config.Addr)
	if err != nil {
		log.Print(err)
		return err
	}
	defer machine.Close()
	log.Printf("connected to %v over IPMI v%v", machine.Address(), machine.Version())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Timeout))
	defer cancel()

	sess, err := machine.NewSession(ctx, &bmc.SessionOpts{
		Username:          config.Username,
		Password:          []byte(config.Password),
		MaxPrivilegeLevel: ipmi.PrivilegeLevelUser,
	})
	if err != nil {
		log.Print(err)
		return err
	}
	defer sess.Close(ctx)
	return task.Execute(ctx, sess)
}

func ExecuteLoop(config config, task Task) {
	ticker := time.NewTicker(time.Duration(config.Ticker) * time.Second)
	defer ticker.Stop()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			return
		default:
		}

		select {
		case <- ticker.C:
			err := Execute(config, task)
			if err != nil {
				log.Print(err)
			}
		case <-signals:
			return
		}
	}
}
