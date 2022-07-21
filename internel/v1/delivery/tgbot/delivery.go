package tgbot

import (
	"github.com/Inari3301/ippinger/internel/v1/usecase"
	tg "gopkg.in/telebot.v3"
)

type Telebot struct {
	u   usecase.UseCase
	bot *tg.Bot
}

func New(u usecase.UseCase, settings tg.Settings) (*Telebot, error) {
	bot, err := tg.NewBot(settings)
	if err != nil {
		return nil, err
	}
	return &Telebot{
		u:   u,
		bot: bot,
	}, nil
}

func (t *Telebot) Start() {
	t.bot.Start()
}
