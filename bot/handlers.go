package bot

import (
	"fmt"
	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"timeBotGolang/config"
)

var (
	WhitelistStateGroup = fsm.NewStateGroup("whitelist")
	WhitelistState      = WhitelistStateGroup.New("startwhitelist")
)

func StartHandler(c tele.Context) error {
	return c.Send("Hello!")
}

func EditWhitelist(c tele.Context, state fsm.Context) error {
	if err := state.Set(WhitelistState); err != nil {
		log.Println(err)
		return c.Send(err)
	}
	return c.Send("Write the IDs separated by space.\nType '0' if you want to set only admin.")
}

func WhitelistStateOnInputIDs(conf *config.Config) func(c tele.Context, state fsm.Context) error {
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
		return c.Send(fmt.Sprintf(
			"Your new whitelist: %v", strings.Trim(fmt.Sprint(conf.Whitelist), "[]")))
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
