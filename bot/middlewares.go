package bot

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"slices"
)

func AccessMiddleware(whitelist []int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !slices.Contains(whitelist, c.Chat().ID) {
				log.Printf("[%d] is not whitelisted.", c.Chat().ID)
				return nil
			}
			return next(c)
		}
	}
}

func AdminAccessMiddleware(id int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().ID != id {
				log.Printf("[%d] is not admin.", c.Sender().ID)
				return nil
			}
			return next(c)
		}
	}
}

func LoggingTestMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			log.Printf("[%d] [%s] executed [%s]\n", c.Sender().ID, c.Sender().Username, c.Text())
			return next(c)
		}
	}
}
