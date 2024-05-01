package bot

import (
	"fmt"
	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func (b *TimeBot) startHandler(c tele.Context) error {
	return b.sendMenu(c)
}

func (b *TimeBot) sendMenu(c tele.Context) error {
	if c.Message().Text == "/start" {
		return c.Send("Choose your option.", menuReplyMarkup)
	}
	return c.Edit("Choose your option.", menuReplyMarkup)
}

func (b *TimeBot) adminMenu(c tele.Context) error {
	const op = "handler.adminMenu"
	if err := c.Send(b.conf.ToString(), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2}); err != nil {
		b.log.Error(op, slog.String("error", err.Error()))
	}
	time.Sleep(time.Second * 1)
	return c.Send("What to change?", adminReplyMarkup)
}

func (b *TimeBot) getTimetable(c tele.Context) error {
	date := b.timetable[0].Format("02-01-2006")

	var timetableMSG = fmt.Sprintf("Timetable for %s:\n", date)
	for _, item := range b.timetable {
		timetableMSG += item.Format("- 15:04") + "\n"
	}
	return c.Edit(timetableMSG, toMenuReplyMarkup)
}

func (b *TimeBot) handleEditWhitelist(c tele.Context, state fsm.Context) error {
	const op = "handler.handleEditWhitelist"
	if err := state.Set(whitelistState); err != nil {
		b.log.Error(op, slog.String("error", err.Error()))
		return c.Send(err)
	}

	message := "Write the IDs separated by space.\nType '0' if you want to set only admin."

	return c.Edit(message)
}

func (b *TimeBot) whitelistStateOnInputIDs(c tele.Context, state fsm.Context) error {
	go func() {
		if err := state.Finish(true); err != nil {
			log.Println(err)
		}
	}()
	idsStr := strings.Fields(c.Text())
	idsInt64, err := strArrToInt64Arr(idsStr)
	if err != nil {
		return c.Send(fmt.Sprintf("Check your input. Error: %s", err))
	}
	b.conf.EditWhitelist(idsInt64)
	return c.Send(b.conf.ToString(),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
}

func (b *TimeBot) handleSetURL(c tele.Context, state fsm.Context) error {
	const op = "handler.handleSetURL"
	if err := state.Set(setURLState); err != nil {
		b.log.Error(op, slog.String("error", err.Error()))
		return c.Send(err)
	}
	return c.Edit("Write the URL for new city.")
}

func (b *TimeBot) setURLStateOnInputURL(c tele.Context, state fsm.Context) error {
	go func() {
		if err := state.Finish(true); err != nil {
			log.Println(err)
		}
	}()
	b.conf.SetCityURL(c.Text())
	return c.Send(b.conf.ToString(),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
}

func (b *TimeBot) handleUpdateTime(c tele.Context, state fsm.Context) error {
	const op = "handler.handleUpdateTime"
	if err := state.Set(setUpdateTimeState); err != nil {
		b.log.Error(op, slog.String("error", err.Error()))
		return c.Send(err)
	}
	return c.Edit("Write the new update time in 15:04 format.")
}

func (b *TimeBot) updateTimeStateOnInputTime(c tele.Context, state fsm.Context) error {
	go func() {
		if err := state.Finish(true); err != nil {
			log.Println(err)
		}
	}()

	_, err := time.Parse("15:04", c.Text())
	if err != nil {
		return c.Send(fmt.Sprintf("Check your input. Error: %s", err))
	}
	b.conf.SetUpdateTime(c.Text())
	return c.Send(b.conf.ToString(),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
}

func (b *TimeBot) handleTimeTill(c tele.Context, state fsm.Context) error {
	const op = "handler.handleTimeTill"
	if err := state.Set(timeTillState); err != nil {
		b.log.Error(op, slog.String("error", err.Error()))
		return c.Send(err)
	}
	return c.Edit("Write the new reminder before in minutes.")
}

func (b *TimeBot) timeTillStateOnInputTime(c tele.Context, state fsm.Context) error {
	go func() {
		if err := state.Finish(true); err != nil {
			log.Println(err)
		}
	}()
	timeTill, err := strconv.Atoi(c.Text())
	if err != nil {
		return c.Send(fmt.Sprintf("Check your input. Error: %s", err))
	}
	b.conf.SetTimeTill(timeTill)
	return c.Send(b.conf.ToString(),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})

}

func strArrToInt64Arr(idsStr []string) ([]int64, error) {
	var idsInt64 []int64
	for _, id := range idsStr {
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		idsInt64 = append(idsInt64, idInt64)
	}
	return idsInt64, nil
}

func (b *TimeBot) getNearTime(c tele.Context) error {
	var nearest time.Time
	for _, t := range b.timetable {
		if t.After(time.Now()) {
			if t.Before(nearest) || nearest.IsZero() {
				nearest = t
			}
		}
	}

	durationTill := time.Until(nearest)
	hours := int(durationTill.Hours())
	minutes := int(durationTill.Minutes()) % 60

	if hours > 0 {
		return c.Send(fmt.Sprintf("Next time in %dh %dm at %v", hours, minutes, nearest.Format("15:04")))
	} else {
		return c.Send(fmt.Sprintf("Next time in %dm at %v", minutes, nearest.Format("15:04")))
	}
}
