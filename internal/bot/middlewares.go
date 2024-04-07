package bot

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"slices"
)

func AccessMiddleware(whitelist *[]int64, log *slog.Logger) tele.MiddlewareFunc {
	const op = "middlewares.AccessMiddleware"
	log.With(slog.String("op", op))

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !slices.Contains(*whitelist, c.Chat().ID) {
				log.Warn("Non whitelisted user tried to use bot",
					slog.Int64("chat_id", c.Chat().ID),
					slog.String("username", fmt.Sprintf("[%v]", c.Chat().Username)))
				return nil
			}
			return next(c)
		}
	}
}

func AdminAccessMiddleware(id int64, log *slog.Logger) tele.MiddlewareFunc {
	const op = "middlewares.AdminAccessMiddleware"
	log.With(slog.String("op", op))
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().ID != id {
				log.Info("Sender is not admin.", slog.Int64("user_id", c.Sender().ID))
				return nil
			}
			return next(c)
		}
	}
}

func LoggingMiddleware(log *slog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Callback() != nil {
				log.Info("Bot received callback",
					slog.Int64("chat_id", c.Sender().ID),
					slog.String("username", c.Sender().Username),
					slog.String("button", c.Callback().Unique))
				return next(c)
			}
			log.Info("Bot received message",
				slog.Int64("chat_id", c.Sender().ID),
				slog.String("username", c.Sender().Username),
				slog.String("message", c.Message().Text))
			return next(c)
		}
	}
}
