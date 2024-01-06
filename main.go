package main

import (
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"time"
	"timeBotGolang/bot"
	"timeBotGolang/parser"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	bot.Whitelist = bot.LoadWhitelist()
	parser.CityURL = os.Getenv("CITY_URL")

	b, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	b.Use(bot.LoggingMiddleware, bot.AccessMiddleware)

	b.Handle("/start", bot.StartHandler)

	b.Start()

}
