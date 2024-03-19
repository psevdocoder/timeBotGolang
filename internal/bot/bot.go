package bot

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
	"timeBotGolang/internal/config"
	"timeBotGolang/internal/scheduler"
)

func InitBot(conf *config.Config) {
	b, err := tele.NewBot(tele.Settings{
		Token:  conf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 60 * time.Second},
	})
	if err != nil {
		log.Println(err)
	}
	log.Printf("Bot %s started\n", b.Me.Username)

	storage := memory.NewStorage()

	InitReplyMarkups()

	b.Use(LoggingMiddleware, AccessMiddleware(conf.Whitelist))
	b.Handle("/start", startHandler)
	b.Handle(&btnGetTimetable, getTimetable)
	b.Handle(&btnToMenu, sendMenu)

	adminOnly := b.Group()
	adminOnly.Use(AdminAccessMiddleware(conf.AdminID))

	adminOnly.Handle("/admin", adminMenu(conf))

	adminOnlyManagerFSM := fsm.NewManager(b, adminOnly, storage, nil)
	adminOnlyManagerFSM.Bind(&btnEditWhitelist, fsm.AnyState, handleEditWhitelist)
	adminOnlyManagerFSM.Bind(tele.OnText, whitelistState, whitelistStateOnInputIDs(conf))
	adminOnlyManagerFSM.Bind(&btnSetURL, fsm.AnyState, handleSetURL)
	adminOnlyManagerFSM.Bind(tele.OnText, setURLState, setURLStateOnInputURL(conf))
	adminOnlyManagerFSM.Bind(&btnUpdateTime, fsm.AnyState, handleUpdateTime)
	adminOnlyManagerFSM.Bind(tele.OnText, setUpdateTimeState, updateTimeStateOnInputTime(conf))
	adminOnlyManagerFSM.Bind(&btnTimeTill, fsm.AnyState, handleTimeTill)
	adminOnlyManagerFSM.Bind(tele.OnText, timeTillState, timeTillStateOnInputTime(conf))

	go scheduler.InitScheduler(conf, b)
	b.Start()
}
