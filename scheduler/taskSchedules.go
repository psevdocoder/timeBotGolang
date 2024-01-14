package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
	"timeBotGolang/config"
	"timeBotGolang/parser"
)

func DailyJobs(s gocron.Scheduler, conf *config.Config, b *tele.Bot) {
	UpdateAt, err := time.Parse("15:04", conf.UpdateTime)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = s.NewJob(gocron.DailyJob(1, gocron.NewAtTimes(
		gocron.NewAtTime(uint(UpdateAt.Hour()), uint(UpdateAt.Minute()), 0))),
		gocron.NewTask(AddTimeTableTasks, s, conf, b))
	if err != nil {
		log.Println(err)
		return
	}
}

func AddTimeTableTasks(s gocron.Scheduler, conf *config.Config, b *tele.Bot) {
	client, _ := parser.NewClient(time.Second * 10)
	timetable := client.GetTimetable(conf.CityURL)
	chatByID, err := b.ChatByID(conf.AdminID)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = b.Send(chatByID, "Requested to site"); err != nil {
		log.Println(err)
		return
	}

	for _, item := range timetable {
		log.Println(item)
		msg := fmt.Sprintf("Next time in %d minutes at %s", conf.TimeTill, item.Format("15:04"))
		notifyAt := item.Add(-conf.TimeTill * time.Minute)
		_, err := s.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Date(
			notifyAt.Year(), notifyAt.Month(), notifyAt.Day(), notifyAt.Hour(), notifyAt.Minute(),
			//time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 5,
			0, 0, time.Local))),
			gocron.NewTask(sendNotification, b, conf, msg))
		if err != nil {
			log.Println(err)
		}
	}
}

func sendNotification(b *tele.Bot, conf *config.Config, msg string) {
	for _, user := range conf.Whitelist {
		chatByID, err := b.ChatByID(user)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err = b.Send(chatByID, msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func InitScheduler(conf *config.Config, b *tele.Bot) {
	s, _ := gocron.NewScheduler()

	AddTimeTableTasks(s, conf, b)
	DailyJobs(s, conf, b)

	s.Start()
}
