package bot

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"timeBotGolang/config"
)

func StartHandler(c tele.Context) error {
	return c.Send("Hello!")
}

func EditWhitelist(conf *config.Config) func(c tele.Context) error {
	return func(c tele.Context) error {
		idsStr := strings.Fields(c.Text())
		log.Println(idsStr)

		newWhitelist, err := StrArrToInt64Arr(idsStr)
		if err != nil {
			return c.Send(fmt.Printf("Check your input. Error: %s", err))
		}
		conf.EditWhitelist(newWhitelist)
		return c.Send("Hello!")

		//TODO FSM

		//variable := []int64{2, 3, 4}
		//
		//conf.EditWhitelist(variable)
		//return c.Send("added")
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
