package tgbot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Handler    func(ctx context.Context, sender *tgbotapi.BotAPI, update tgbotapi.Update) context.Context
	Middleware func(ctx context.Context, update tgbotapi.Update) context.Context

	Sender interface {
		Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
		GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error)
	}

	Router interface {
		Handler(pattern string, handler Handler)
		Middleware(handler Middleware)
		Match(pattern string) (Handler, bool)
	}
)
