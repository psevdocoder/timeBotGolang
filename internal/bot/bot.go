package bot

import (
	"errors"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"time"
	"timeBotGolang/internal/config"
	"timeBotGolang/internal/parser"
)

type TimeBot struct {
	Bot       *tele.Bot
	conf      *config.Config
	scheduler gocron.Scheduler
	timetable []time.Time
	log       *slog.Logger
}

func NewBot(log *slog.Logger, conf *config.Config) *TimeBot {

	b, err := tele.NewBot(tele.Settings{
		Token:  conf.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 60 * time.Second},
	})
	if err != nil {
		log.Error("Failed to start bot", slog.String("error", err.Error()))
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Error("Failed to create scheduler", slog.String("error", err.Error()))
	}

	return &TimeBot{Bot: b, conf: conf, scheduler: s, log: log}
}

func (b *TimeBot) Start() {
	const op = "bot.Start"
	log := b.log.With(slog.String("op", op))

	storage := memory.NewStorage()

	InitReplyMarkups()

	b.Bot.Use(LoggingMiddleware(b.log), AccessMiddleware(&b.conf.Whitelist, b.log))
	b.Bot.Handle("/start", b.startHandler)
	b.Bot.Handle(&btnGetTimetable, b.getTimetable)
	b.Bot.Handle(&btnToMenu, b.sendMenu)
	b.Bot.Handle("/near", b.getNearTime)

	adminOnly := b.Bot.Group()
	adminOnly.Use(AdminAccessMiddleware(b.conf.AdminID, b.log))

	adminOnly.Handle("/admin", b.adminMenu)

	adminOnlyManagerFSM := fsm.NewManager(b.Bot, adminOnly, storage, nil)
	adminOnlyManagerFSM.Bind(&btnEditWhitelist, fsm.AnyState, b.handleEditWhitelist)
	adminOnlyManagerFSM.Bind(tele.OnText, whitelistState, b.whitelistStateOnInputIDs)
	adminOnlyManagerFSM.Bind(&btnSetURL, fsm.AnyState, b.handleSetURL)
	adminOnlyManagerFSM.Bind(tele.OnText, setURLState, b.setURLStateOnInputURL)
	adminOnlyManagerFSM.Bind(&btnUpdateTime, fsm.AnyState, b.handleUpdateTime)
	adminOnlyManagerFSM.Bind(tele.OnText, setUpdateTimeState, b.updateTimeStateOnInputTime)
	adminOnlyManagerFSM.Bind(&btnTimeTill, fsm.AnyState, b.handleTimeTill)
	adminOnlyManagerFSM.Bind(tele.OnText, timeTillState, b.timeTillStateOnInputTime)

	b.AddTimeTableTasks()
	b.DailyJobs()
	go b.scheduler.Start()
	go b.Bot.Start()
	log.Info("Bot started", slog.String("bot_name", b.Bot.Me.Username))
}

func (b *TimeBot) DailyJobs() {
	const op = "bot.DailyJobs"
	log := b.log.With(slog.String("op", op))

	b.log.With(slog.String("op", op))
	UpdateAt, err := time.Parse("15:04", b.conf.UpdateTime)
	if err != nil {
		b.log.Error("Failed to parse update time", slog.String("error", err.Error()))
		return
	}

	_, err = b.scheduler.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(
		gocron.NewAtTime(uint(UpdateAt.Hour()), uint(UpdateAt.Minute()), 0))),
		gocron.NewTask(b.AddTimeTableTasks))
	if err != nil {
		log.Error("Failed to add job", slog.String("error", err.Error()))
		return
	}
}

func (b *TimeBot) AddTimeTableTasks() {
	const op = "bot.AddTimeTableTasks"
	log := b.log.With(slog.String("op", op))

	client, err := parser.NewClient(time.Second*10, b.log)
	if err != nil {
		log.Error("Failed to create client", slog.String("error", err.Error()))
		return
	}

	b.timetable = client.GetTimetable(b.conf.CityURL)
	chatByID, err := b.Bot.ChatByID(b.conf.AdminID)
	if err != nil {
		log.Error("Failed to load chat id", slog.String("error", err.Error()))
		return
	}

	if _, err = b.Bot.Send(chatByID, "Requested to site"); err != nil {
		log.Error("Failed to send message", slog.String("error", err.Error()))
		return
	}

	log.Info("Timetable loaded", slog.String("timetable", fmt.Sprint(b.timetable)))
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
			if errors.Is(err, gocron.ErrOneTimeJobStartDateTimePast) {
				log.Debug("Past time was skipped", slog.String("error", err.Error()))
				continue
			}
			log.Error("Failed to add job", slog.String("error", err.Error()))
		}
	}

}

func (b *TimeBot) sendNotification(msg string) {
	const op = "bot.sendNotification"
	log := b.log.With(slog.String("op", op))

	log.Info("Sending notification", slog.String("message", msg))
	for _, user := range b.conf.Whitelist {
		chatByID, err := b.Bot.ChatByID(user)
		if err != nil {
			log.Error("Failed to load chat id", slog.String("error", err.Error()))
			return
		}
		if _, err = b.Bot.Send(chatByID, msg, menuReplyMarkup); err != nil {
			log.Error("Failed to send message", slog.String("error", err.Error()))
			return
		}
	}
}

func (b *TimeBot) Stop() {
	const op = "bot.Stop"
	log := b.log.With(slog.String("op", op))

	log.Info("Stopping bot...")

	adminChat, err := b.Bot.ChatByID(b.conf.AdminID)
	if err != nil {
		log.Error("Failed to load chat id", slog.String("error", err.Error()))
		return
	}

	_, _ = b.Bot.Send(adminChat, "Bot stopped")
	b.Bot.Stop()
}
