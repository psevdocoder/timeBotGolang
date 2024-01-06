package bot

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

var Whitelist []int64

func LoadWhitelist() []int64 {
	whitelistStr := os.Getenv("WHITE_LIST")
	whiteListInt64 := func(s string) []int64 {
		var result []int64
		for _, str := range strings.Split(s, ",") {
			if num, err := strconv.ParseInt(str, 10, 64); err == nil {
				result = append(result, num)
			}
		}
		return result
	}(whitelistStr)
	return whiteListInt64
}

func AccessMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !slices.Contains(Whitelist, c.Chat().ID) {
			log.Printf("[%d] is not whitelisted.", c.Chat().ID)
			return nil
		}
		return next(c)
	}
}

func LoggingMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		log.Printf("[%d] executed [%s]", c.Sender().ID, c.Text())
		return next(c)
	}
}
