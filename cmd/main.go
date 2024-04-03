package main

import (
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"timeBotGolang/internal/bot"
	"timeBotGolang/internal/config"
)

func main() {
	log, file := SetupLogger()
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	log.Info("Starting bot...")

	conf, err := config.NewConfig()
	if err != nil {
		log.Error("Failed to load config", slog.String("error", err.Error()))
	}

	log.Debug("Current config", slog.Any("config", conf))

	myBot := bot.NewBot(log, conf)

	myBot.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	myBot.Stop()
}

func SetupLogger() (*slog.Logger, *os.File) {
	var logger *slog.Logger

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return logger, file
}
