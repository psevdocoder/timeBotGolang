package bot

import tele "gopkg.in/telebot.v3"

var (
	menuReplyMarkup = &tele.ReplyMarkup{}
	btnGetTimetable = menuReplyMarkup.Data("Timetable", "timetable")

	adminReplyMarkup = &tele.ReplyMarkup{}
	btnEditWhitelist = adminReplyMarkup.Data("Whitelist", "whitelist")
	btnSetURL        = adminReplyMarkup.Data("Set URL", "setUrl")
	btnUpdateTime    = adminReplyMarkup.Data("Update Time", "updateTime")
	btnTimeTill      = adminReplyMarkup.Data("Time Till", "timeTill")

	toMenuReplyMarkup = &tele.ReplyMarkup{}
	btnToMenu         = toMenuReplyMarkup.Data("Back", "back")
)

func InitReplyMarkups() {
	menuReplyMarkup.Inline(
		menuReplyMarkup.Row(btnGetTimetable),
	)

	adminReplyMarkup.Inline(
		adminReplyMarkup.Row(btnEditWhitelist, btnSetURL),
		adminReplyMarkup.Row(btnUpdateTime, btnTimeTill),
		adminReplyMarkup.Row(btnToMenu),
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
