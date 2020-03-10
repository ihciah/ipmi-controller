package ipmi_controller

import (
	"errors"
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TelegramBot struct {
	*IPMI
	bot *tb.Bot
}

func NewTelegramBot(ipmi *IPMI) (*TelegramBot, error) {
	if ipmi.TelegramConfig.Token == "" {
		return nil, errors.New("can not load telegram config")
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  ipmi.TelegramConfig.Token,
		URL:    ipmi.TelegramConfig.URL,
		Poller: &tb.LongPoller{Timeout: time.Duration(ipmi.TelegramConfig.PollTimeout) * time.Second},
	})

	return &TelegramBot{
		IPMI: ipmi,
		bot:  bot,
	}, err
}

func (t *TelegramBot) Serve() {
	getStatusBtn := tb.ReplyButton{Text: "GetStatus"}
	getTemperatureBtn := tb.ReplyButton{Text: "GetTemperature"}
	setPowerOnBtn := tb.ReplyButton{Text: "SetPowerOn"}
	setPowerOffBtn := tb.ReplyButton{Text: "SetPowerOff"}

	replyKeys := [][]tb.ReplyButton{
		{getStatusBtn, getTemperatureBtn},
		{setPowerOnBtn, setPowerOffBtn},
	}

	t.bot.Handle(&getStatusBtn, t.callbackFactory(t.GetStatus))
	t.bot.Handle(&getTemperatureBtn, t.callbackFactory(t.GetTemperature))
	t.bot.Handle(&setPowerOnBtn, t.callbackFactory(t.SetPowerOn))
	t.bot.Handle(&setPowerOffBtn, t.callbackFactory(t.SetPowerOff))
	t.bot.Handle("/ipmi", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		t.bot.Send(m.Sender, "Please choose the button", &tb.ReplyMarkup{
			ReplyKeyboard: replyKeys,
		})
	})
	log.Print("telegram bot started")
	t.bot.Start()
}

func (t *TelegramBot) callbackFactory(f func() (string, error)) func(m *tb.Message) {
	fun := f
	return func(m *tb.Message) {
		authorized := false
		for _, id := range t.TelegramConfig.Admin {
			if id == m.Sender.ID {
				authorized = true
				break
			}
		}
		log.Print(m.Sender.ID, authorized)
		if !authorized {
			t.bot.Reply(m, fmt.Sprintf("Unauthorized user %d.", m.Sender.ID))
			return
		}
		msg, err := t.bot.Reply(m, "Command executing...")
		if err != nil {
			return
		}
		output, err := fun()
		if err != nil {
			output = err.Error()
		}
		t.bot.Edit(msg, output)
	}
}
