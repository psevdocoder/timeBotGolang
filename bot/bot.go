package bot

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
	"timeBotGolang/config"
	"timeBotGolang/scheduler"
)

func InitBot(conf *config.Config) {
	b, err := tele.NewBot(tele.Settings{
		Token:  conf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bot %s started\n", b.Me.Username)

	storage := memory.NewStorage()

	InitReplyMarkups()

	b.Use(LoggingMiddleware, AccessMiddleware(conf.Whitelist))
	b.Handle("/start", StartHandler)
	b.Handle(&btnToMenu, SendMenu)
	b.Handle(&btnGetTimetable, GetTimetable)

	adminOnly := b.Group()
	adminOnly.Use(AdminAccessMiddleware(conf.AdminID))
	adminOnlyManager := fsm.NewManager(b, adminOnly, storage, nil)
	adminOnlyManager.Bind(&btnEditWhitelist, fsm.AnyState, handleEditWhitelist)
	adminOnlyManager.Bind(tele.OnText, whitelistState, whitelistStateOnInputIDs(conf))
	adminOnlyManager.Bind(&btnSetURL, fsm.AnyState, handleSetURL)
	adminOnlyManager.Bind(tele.OnText, setURLState, setURLStateOnInputURL(conf))
	adminOnlyManager.Bind(&btnUpdateTime, fsm.AnyState, handleUpdateTime)
	adminOnlyManager.Bind(tele.OnText, setUpdateTimeState, updateTimeStateOnInputTime(conf))
	adminOnlyManager.Bind(&btnTimeTill, fsm.AnyState, handleTimeTill)
	adminOnlyManager.Bind(tele.OnText, timeTillState, timeTillStateOnInputTime(conf))

	go scheduler.InitScheduler(conf, b)
	b.Start()
}
