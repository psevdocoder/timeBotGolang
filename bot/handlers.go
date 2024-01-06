package bot

import tele "gopkg.in/telebot.v3"

func StartHandler(c tele.Context) error {
	return c.Send("Hello!")
}
