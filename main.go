package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"timeBotGolang/bot"
	"timeBotGolang/config"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	go bot.InitBot(conf)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Stopping bot...")
}
