package tgmux

import (
	"context"
	"github.com/Inari3301/ippinger/internel/v1/delivery/tgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	middlewares []tgbot.Middleware
	handlers    map[string]tgbot.Handler
}

func New() *Router {
	return &Router{
		middlewares: make([]tgbot.Middleware, 0),
		handlers:    make(map[string]tgbot.Handler),
	}
}

func (r *Router) Handler(pattern string, handler tgbot.Handler) {
	r.handlers[pattern] = handler
}

func (r *Router) Middleware(middleware tgbot.Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Match(pattern string) (tgbot.Handler, bool) {
	handler, exist := r.handlers[pattern]
	if !exist {
		return nil, false
	}

	return func(ctx context.Context, sender *tgbotapi.BotAPI, update tgbotapi.Update) context.Context {
		for _, middleware := range r.middlewares {
			middleware(ctx, update)
		}
		return handler(ctx, sender, update)
	}, true
}
