package tgbot

import (
	"context"
	"github.com/Inari3301/ippinger/internel/v1/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

var (
	commands = strings.Join([]string{"/start", "/ping", "/ping_csv"}, "\n")
)

const (
	enter State = 1
)

type Processor struct {
	U usecase.UseCase
}

func (p Processor) Start(ctx context.Context, sender Sender, update tgbotapi.Update) context.Context {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, commands)
	_, _ = sender.Send(msg)
	return ctx
}

func (p Processor) Ping(ctx context.Context, sender Sender, update tgbotapi.Update) context.Context {
	s := ctx.Value(state).(State)
	switch s {
	case None:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "enter ip")
		ctx = context.WithValue(ctx, state, enter)
		ctx = context.WithValue(ctx, currentHandler, p.Ping)
		_, _ = sender.Send(msg)
	case enter:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "enadasdasasd")
		ctx = context.WithValue(ctx, state, None)
		_, _ = sender.Send(msg)
	}
	return ctx
}
