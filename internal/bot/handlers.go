package bot

import (
	"fmt"
	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"time"
	"timeBotGolang/internal/config"
	"timeBotGolang/internal/scheduler"
)

func startHandler(c tele.Context) error {
	return sendMenu(c)
}

func sendMenu(c tele.Context) error {
	if c.Message().Text == "/start" {
		return c.Send("Choose your option.", menuReplyMarkup)
	}
	return c.Edit("Choose your option.", menuReplyMarkup)
}

func adminMenu(conf *config.Config) func(c tele.Context) error {
	return func(c tele.Context) error {
		if err := c.Send(conf.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2}); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 1)
		return c.Send("What to change?", adminReplyMarkup)
	}
}

func getTimetable(c tele.Context) error {
	date := scheduler.Timetable[0].Format("02-01-2006")

	var timetableMSG = fmt.Sprintf("Timetable for %s:\n", date)
	for _, item := range scheduler.Timetable {
		timetableMSG += item.Format("- 15:04") + "\n"
	}
	return c.Edit(timetableMSG, toMenuReplyMarkup)
}

func handleEditWhitelist(c tele.Context, state fsm.Context) error {
	if err := state.Set(whitelistState); err != nil {
		log.Println(err)
		return c.Send(err)
	}
	return c.Edit("Write the IDs separated by space.\nType '0' if you want to set only admin.")
}

func whitelistStateOnInputIDs(conf *config.Config) func(c tele.Context, state fsm.Context) error {
	return func(c tele.Context, state fsm.Context) error {
		go func() {
			if err := state.Finish(true); err != nil {
				log.Println(err)
			}
		}()

		idsStr := strings.Fields(c.Text())
		idsInt64, err := StrArrToInt64Arr(idsStr)
		if err != nil {
			return c.Send(fmt.Sprintf("Check your input. Error: %s", err))
		}
		conf.EditWhitelist(idsInt64)
		return c.Send(conf.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
	}
}

func handleSetURL(c tele.Context, state fsm.Context) error {
	if err := state.Set(setURLState); err != nil {
		log.Println(err)
		return c.Send(err)
	}
	return c.Edit("Write the URL for new city.")
}

func setURLStateOnInputURL(conf *config.Config) func(c tele.Context, state fsm.Context) error {
	return func(c tele.Context, state fsm.Context) error {
		go func() {
			if err := state.Finish(true); err != nil {
				log.Println(err)
			}
		}()
		conf.SetCityURL(c.Text())
		return c.Send(conf.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
	}
}

func handleUpdateTime(c tele.Context, state fsm.Context) error {
	if err := state.Set(setUpdateTimeState); err != nil {
		log.Println(err)
		return c.Send(err)
	}
	return c.Edit("Write the new update time in 15:04 format.")
}

func updateTimeStateOnInputTime(conf *config.Config) func(c tele.Context, state fsm.Context) error {
	return func(c tele.Context, state fsm.Context) error {
		go func() {
			if err := state.Finish(true); err != nil {
				log.Println(err)
			}
		}()

		_, err := time.Parse("15:04", c.Text())
		if err != nil {
			return c.Send(fmt.Sprintf("Check your input. Error: %s", err))
		}
		conf.SetUpdateTime(c.Text())
		return c.Send(conf.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
	}
}

func handleTimeTill(c tele.Context, state fsm.Context) error {
	if err := state.Set(timeTillState); err != nil {
		log.Println(err)
		return c.Send(err)
	}
	return c.Edit("Write the new reminder before in minutes.")
}

func timeTillStateOnInputTime(conf *config.Config) func(c tele.Context, state fsm.Context) error {
	return func(c tele.Context, state fsm.Context) error {
		go func() {
			if err := state.Finish(true); err != nil {
				log.Println(err)
			}
		}()
		timeTill, err := strconv.Atoi(c.Text())
		if err != nil {
			return c.Send(fmt.Sprintf("Check your input. Error: %s", err))
		}
		conf.SetTimeTill(timeTill)
		return c.Send(conf.String(), &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, ReplyMarkup: adminReplyMarkup})
	}
}

func StrArrToInt64Arr(idsStr []string) ([]int64, error) {
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
