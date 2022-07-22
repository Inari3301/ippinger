package main

import (
	"github.com/Inari3301/ippinger/internel/v1/delivery/tgbot"
	"github.com/Inari3301/ippinger/internel/v1/delivery/tgbot/tgmux"
	"github.com/Inari3301/ippinger/internel/v1/store/memstore"
	"github.com/Inari3301/ippinger/internel/v1/usecase/pkgusecase"
	"log"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	s, err := memstore.New(memstore.Options{
		Path:         os.Args[1],
		DumpInterval: 15,
		BatchSize:    100,
	})

	if err != nil {
		log.Fatalln(err)
	}

	u := pkgusecase.New(s)
	router := tgmux.New()
	proc := tgbot.Processor{
		U: u,
	}
	router.Handler("/start", proc.Start)
	router.Handler("/ping", proc.Ping)

	bot, err := tgbot.New(tgbot.Options{
		Token: "5408879578:AAGXUy245KzdSC9fyXBAJ6StXUYOsJhdhwE",
	}, router)

	if err != nil {
		log.Fatalln(err)
	}

	bot.Start()
}
