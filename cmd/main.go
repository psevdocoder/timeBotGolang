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

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

func main() {
	conf, err := config.NewConfig()

	log, file := SetupLogger(conf.Env)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	log.Info("Starting bot...")

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

func SetupLogger(env string) (*slog.Logger, *os.File) {
	var logger *slog.Logger

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	switch env {
	case EnvDev:
		logger = slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger, file
}
