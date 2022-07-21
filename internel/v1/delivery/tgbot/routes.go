package tgbot

import tg "gopkg.in/telebot.v3"

func (t *Telebot) newRoutes() {
	t.bot.Handle("start", func(c tg.Context) error {
		return c.Send("Hello")
	})
}
