package tgbot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Inari3301/ippinger/internel/v1/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
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
		ctx = context.WithValue(ctx, currentHandler, "/ping")
		_, _ = sender.Send(msg)
	case enter:
		ctx = context.WithValue(ctx, state, None)
		str, err := p.U.Ping(update.Message.Text, 10*time.Second)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			_, _ = sender.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, %d", str.IP, str.Duration))
			_, _ = sender.Send(msg)
		}
	}
	return ctx
}

func (p Processor) PingCsv(ctx context.Context, sender Sender, update tgbotapi.Update) context.Context {
	s := ctx.Value(state).(State)
	switch s {
	case None:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "enter ips")
		ctx = context.WithValue(ctx, state, enter)
		ctx = context.WithValue(ctx, currentHandler, "/ping_csv")
		_, _ = sender.Send(msg)
	case enter:
		ctx = context.WithValue(ctx, state, None)
		str, err := p.U.PingByCsv([]byte(update.Message.Text))
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			_, _ = sender.Send(msg)
		} else {
			b, _ := json.Marshal(str.PingResults)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(b))
			_, _ = sender.Send(msg)
		}
	}
	return ctx
}
