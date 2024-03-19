package bot

import "github.com/vitaliy-ukiru/fsm-telebot"

var (
	whitelistStateGroup = fsm.NewStateGroup("whitelist")
	whitelistState      = whitelistStateGroup.New("whitelistState")

	setURLStateGroup = fsm.NewStateGroup("setURL")
	setURLState      = setURLStateGroup.New("setURLState")

	setUpdateTimeStateGroup = fsm.NewStateGroup("setUpdateTime")
	setUpdateTimeState      = setUpdateTimeStateGroup.New("setUpdateTimeState")

	timeTillStateGroup = fsm.NewStateGroup("timeTill")
	timeTillState      = timeTillStateGroup.New("timeTillState")
)
