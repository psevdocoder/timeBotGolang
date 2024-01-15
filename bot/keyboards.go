package bot

import tele "gopkg.in/telebot.v3"

var (
	menuReplyMarkup  = &tele.ReplyMarkup{}
	btnGetTimetable  = menuReplyMarkup.Data("Timetable", "timetable")
	btnEditWhitelist = menuReplyMarkup.Data("Whitelist", "whitelist")
	btnSetURL        = menuReplyMarkup.Data("Set URL", "setUrl")
	btnUpdateTime    = menuReplyMarkup.Data("Update Time", "updateTime")
	btnTimeTill      = menuReplyMarkup.Data("Time Till", "timeTill")

	toMenuReplyMarkup = &tele.ReplyMarkup{}
	btnToMenu         = toMenuReplyMarkup.Data("Back", "back")
)

func InitReplyMarkups() {
	menuReplyMarkup.Inline(
		menuReplyMarkup.Row(btnGetTimetable),
		menuReplyMarkup.Row(btnEditWhitelist, btnSetURL),
		menuReplyMarkup.Row(btnUpdateTime, btnTimeTill),
	)

	toMenuReplyMarkup.Inline(
		toMenuReplyMarkup.Row(btnToMenu),
	)
}

//func qwe() {
//	menu.Reply(
//		menu.Row(btnHelp),
//		menu.Row(btnSettings),
//	)
//	menuReplyMarkup.Inline(
//		menuReplyMarkup.Row(btnGetTimetable, btnNext),
//	)
//
//	b.Handle("/start", func(c tele.Context) error {
//		return c.Send("Hello!", menu)
//	})
//
//	// On reply button pressed (message)
//	b.Handle(&btnHelp, func(c tele.Context) error {
//		return c.Edit("Here is some help: ...")
//	})
//
//	// On inline button pressed (callback)
//	b.Handle(&btnGetTimetable, func(c tele.Context) error {
//		return c.Respond()
//	})
//}
