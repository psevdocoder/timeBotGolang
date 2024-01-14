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

	b.Use(LoggingTestMiddleware(), AccessMiddleware(conf.Whitelist))
	b.Handle("/start", StartHandler)

	adminOnly := b.Group()
	adminOnly.Use(AdminAccessMiddleware(conf.AdminID))
	adminOnlyManager := fsm.NewManager(b, adminOnly, storage, nil)
	adminOnlyManager.Bind("/whitelist", fsm.AnyState, EditWhitelist)
	adminOnlyManager.Bind(tele.OnText, WhitelistState, WhitelistStateOnInputIDs(conf))

	go scheduler.InitScheduler(conf, b)
	b.Start()
}
