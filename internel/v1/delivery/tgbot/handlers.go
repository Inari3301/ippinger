package tgbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Inari3301/ippinger/internel/v1/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
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

func (p Processor) Start(ctx context.Context, sender *tgbotapi.BotAPI, update tgbotapi.Update) context.Context {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, commands)
	_, _ = sender.Send(msg)
	return ctx
}

func (p Processor) Ping(ctx context.Context, sender *tgbotapi.BotAPI, update tgbotapi.Update) context.Context {
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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s, %v"+"\n"+"%s", str.IP, str.Duration, commands))
			_, _ = sender.Send(msg)
		}
	}
	return ctx
}

func (p Processor) PingCsv(ctx context.Context, sender *tgbotapi.BotAPI, update tgbotapi.Update) context.Context {
	s := ctx.Value(state).(State)
	switch s {
	case None:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "enter ips")
		ctx = context.WithValue(ctx, state, enter)
		ctx = context.WithValue(ctx, currentHandler, "/ping_csv")
		_, _ = sender.Send(msg)
	case enter:
		ctx = context.WithValue(ctx, state, None)
		f, err := sender.GetFile(tgbotapi.FileConfig{
			FileID: update.Message.Document.FileID,
		})
		if err != nil {
			_, _ = sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		r, _ := http.NewRequest("GET", fmt.Sprintf(tgbotapi.FileEndpoint, sender.Token, f.FilePath), &bytes.Reader{})
		resp, err := sender.Client.Do(r)
		if err != nil {
			_, _ = sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		b := make([]byte, 1000)
		defer resp.Body.Close()
		n, err := resp.Body.Read(b)
		fmt.Println(string(b))
		if err != nil {
			_, _ = sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}
		b = b[:n]
		str, err := p.U.PingByCsv(b)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			_, _ = sender.Send(msg)
		}
		b, err = json.Marshal(str)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(b))
		_, _ = sender.Send(msg)
	}
	return ctx
}
