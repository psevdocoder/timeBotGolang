package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"timeBotGolang/internal/bot"
	"timeBotGolang/internal/config"
)

func main() {
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file:", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	log.SetOutput(io.MultiWriter(file, os.Stdout))

	conf, err := config.NewConfig()
	if err != nil {
		log.Println(err)
	}

	go bot.InitBot(conf)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Stopping bot...")
}
