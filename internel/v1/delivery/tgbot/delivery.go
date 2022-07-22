package tgbot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type (
	ContextKeyType uint64
	State          uint64
)

const (
	state          ContextKeyType = 1
	currentHandler ContextKeyType = 2
)

const (
	None State = iota
)

type Options struct {
	Token string
}

type Telebot struct {
	bot     *tgbotapi.BotAPI
	router  Router
	lock    sync.Mutex
	userCtx map[int64]context.Context
}

func New(opt Options, router Router) (*Telebot, error) {
	bot, err := tgbotapi.NewBotAPI(opt.Token)
	bot.Debug = true
	if err != nil {
		return nil, err
	}
	return &Telebot{
		bot:     bot,
		router:  router,
		userCtx: make(map[int64]context.Context),
	}, nil
}

func (t *Telebot) Start() {
	u := tgbotapi.NewUpdate(0)
	c := t.bot.GetUpdatesChan(u)
	for update := range c {
		t.lock.Lock()
		ctx, exist := t.userCtx[update.Message.From.ID]
		if !exist {
			ctx = context.Background()
			ctx = context.WithValue(ctx, state, None)
			t.userCtx[update.Message.From.ID] = ctx
		}
		t.lock.Unlock()

		var handler Handler
		if ctx.Value(state).(State) == None {
			handler, exist = t.router.Match(update.Message.Text)
			if !exist {
				t.badRequest(update)
				continue
			}
		} else {
			handler, exist = ctx.Value(currentHandler).(Handler)
			if !exist {
				log.Println("handler does not exists")
				t.badRequest(update)
				continue
			}
		}

		go func(u tgbotapi.Update) {
			cont := handler(ctx, t.bot, u)
			t.lock.Lock()
			t.userCtx[u.Message.From.ID] = cont
			t.lock.Unlock()
		}(update)
	}
}

func (t *Telebot) badRequest(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/start")
	m, err := t.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(m)
}
