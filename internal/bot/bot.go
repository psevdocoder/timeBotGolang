package bot

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
	"timeBotGolang/internal/config"
	"timeBotGolang/internal/parser"
)

type TimeBot struct {
	Bot       *tele.Bot
	conf      *config.Config
	scheduler gocron.Scheduler
	timetable []time.Time
}

func NewBot(conf *config.Config) *TimeBot {

	b, err := tele.NewBot(tele.Settings{
		Token:  conf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 60 * time.Second},
	})
	if err != nil {
		log.Println(err)
	}

	s, _ := gocron.NewScheduler()

	return &TimeBot{Bot: b, conf: conf, scheduler: s}
}

func (b *TimeBot) Start() {

	storage := memory.NewStorage()

	InitReplyMarkups()

	b.Bot.Use(LoggingMiddleware, AccessMiddleware(&b.conf.Whitelist))
	b.Bot.Handle("/start", startHandler)
	b.Bot.Handle(&btnGetTimetable, getTimetable(b))
	b.Bot.Handle(&btnToMenu, sendMenu)

	adminOnly := b.Bot.Group()
	adminOnly.Use(AdminAccessMiddleware(b.conf.AdminID))

	adminOnly.Handle("/admin", adminMenu(b.conf))

	adminOnlyManagerFSM := fsm.NewManager(b.Bot, adminOnly, storage, nil)
	adminOnlyManagerFSM.Bind(&btnEditWhitelist, fsm.AnyState, handleEditWhitelist)
	adminOnlyManagerFSM.Bind(tele.OnText, whitelistState, whitelistStateOnInputIDs(b.conf))
	adminOnlyManagerFSM.Bind(&btnSetURL, fsm.AnyState, handleSetURL)
	adminOnlyManagerFSM.Bind(tele.OnText, setURLState, setURLStateOnInputURL(b.conf))
	adminOnlyManagerFSM.Bind(&btnUpdateTime, fsm.AnyState, handleUpdateTime)
	adminOnlyManagerFSM.Bind(tele.OnText, setUpdateTimeState, updateTimeStateOnInputTime(b.conf))
	adminOnlyManagerFSM.Bind(&btnTimeTill, fsm.AnyState, handleTimeTill)
	adminOnlyManagerFSM.Bind(tele.OnText, timeTillState, timeTillStateOnInputTime(b.conf))

	b.AddTimeTableTasks()
	b.DailyJobs()
	go b.scheduler.Start()
	go b.Bot.Start()
	log.Printf("Bot %s started\n", b.Bot.Me.Username)
}

func (b *TimeBot) DailyJobs() {
	UpdateAt, err := time.Parse("15:04", b.conf.UpdateTime)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = b.scheduler.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(
		gocron.NewAtTime(uint(UpdateAt.Hour()), uint(UpdateAt.Minute()), 0))),
		gocron.NewTask(b.AddTimeTableTasks))
	if err != nil {
		log.Println(err)
		return
	}
}

func (b *TimeBot) AddTimeTableTasks() {
	client, _ := parser.NewClient(time.Second * 10)
	b.timetable = client.GetTimetable(b.conf.CityURL)
	chatByID, err := b.Bot.ChatByID(b.conf.AdminID)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = b.Bot.Send(chatByID, "Requested to site"); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Loaded timetable: %v", b.timetable)
	for _, item := range b.timetable {
		msg := fmt.Sprintf("Next time in %d minutes at %s", b.conf.TimeTill, item.Format("15:04"))

		minusTime := time.Duration(b.conf.TimeTill) * time.Minute
		notifyAt := item.Add(-minusTime)
		_, err = b.scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Date(
			notifyAt.Year(), notifyAt.Month(), notifyAt.Day(), notifyAt.Hour(), notifyAt.Minute(),
			//time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 5,
			0, 0, time.Local))),
			gocron.NewTask(b.sendNotification, msg))
		if err != nil {
			log.Println(err)
		}
	}

}

func (b *TimeBot) sendNotification(msg string) {
	log.Printf("Sending notification [%s]", msg)
	for _, user := range b.conf.Whitelist {
		chatByID, err := b.Bot.ChatByID(user)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err = b.Bot.Send(chatByID, msg, menuReplyMarkup); err != nil {
			log.Println(err)
			return
		}
	}
}
