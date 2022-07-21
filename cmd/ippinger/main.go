package main

import (
	"github.com/Inari3301/ippinger/internel/v1/delivery/tgbot"
	"github.com/Inari3301/ippinger/internel/v1/store/memstore"
	"github.com/Inari3301/ippinger/internel/v1/usecase/pkgusecase"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	s, err := memstore.New(memstore.Options{
		Path:         os.Args[1],
		DumpInterval: 15,
		BatchSize:    100,
	})

	if err != nil {
		log.Fatalln(err)
	}

	runtime.LockOSThread()
	u := pkgusecase.New(s)
	bot, err := tgbot.New(u, telebot.Settings{
		Token:  os.Getenv("IPPINGER_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalln(err)
	}

	bot.Start()
}
