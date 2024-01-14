# TimeBot Golang
## Description

TimeBot is a Telegram bot that allows you to get the notifications in some cities in advance.  
It makes requests to the certain site, which has structure similar to the one below:

| Month | Column1 | Column2 | Column3 |
| ---- | :--- | :--- | :--- |
| 1 (пн) | 05:30 | 08:00 | 12:00 |
| 2 (вт) | 06:15 | 08:30 | 11:45 |
| 3 (ср) | 06:00 | 08:15 | 12:30 |
| 4 (чт) | 06:20 | 08:40 | 12:10 |
| 5 (пт) | 06:35 | 08:50 | 12:25 |
| 6 (сб) | 06:10 | 08:25 | 12:20 |
| 7 (вс) | 06:25 | 08:45 | 12:15 |
| 8 (пн) | 06:05 | 08:35 | 12:05 |
| 9 (вт) | 06:30 | 08:55 | 12:35 |
| 10 (ср) | 06:50 | 08:20 | 12:50 |
| 11 (чт) | 06:40 | 08:10 | 12:30 |

## Features:
- Users' whitelist by telegram ID.
- Using the Admin's telegram ID for administration commands.
- In bot config editor (For admin).
- Time reminder a certain number of minutes before.
- Daily updating of timetable for current day.

## Used dependencies
### Internal
- net/http for requests with http.Client
- encoding/json for working with config.json
- slices
- strconv, strings
- log
- os
- fmt
- time

### External
- [Telebot](https://github.com/tucnak/telebot) - golang framework for telegram bots development.
- [Telebot-FSM](https://github.com/vitaliy-ukiru/fsm-telebot) - finite state machine for `Telebot`. Used for step-by-step dialogs realization.
- [Gocron](https://github.com/go-co-op/gocron) - package for jobs scheduling. Used for sending reminding messages in advance by cron.
- [Goquery](https://github.com/PuerkitoBio/goquery) - provides parsing tools. Used for getting useful information from http response.
  