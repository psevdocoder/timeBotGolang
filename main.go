package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
	"timeBotGolang/bot"
	"timeBotGolang/config"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(tele.Settings{
		Token:  conf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}
	b.Use(bot.LoggingTestMiddleware(), bot.AccessMiddleware(conf.Whitelist))
	b.Handle("/start", bot.StartHandler)

	b.Handle("/whitelist", bot.EditWhitelist(conf))

	b.Start()

}
