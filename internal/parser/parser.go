package parser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type roundTripper struct {
	next http.RoundTripper
	log  *slog.Logger
}

type Client struct {
	*http.Client
	log *slog.Logger
}

func NewClient(timeout time.Duration, log *slog.Logger) (*Client, error) {
	const op = "parser.NewClient"
	log = log.With(slog.String("op", op))
	log.Debug("Creating client", slog.Duration("timeout", timeout))

	if timeout == 0 {
		return nil, errors.New("timeout must be greater than 0")
	}
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
			Transport: &roundTripper{
				next: http.DefaultTransport,
				log:  log,
			},
		},
		log: log,
	}, nil
}

func (t *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	const op = "parser.roundTripper.RoundTrip"
	log := t.log.With(slog.String("op", op))

	req.Header.Set("User-Agent", userAgent)

	log.Debug("Sending request", slog.String("url", fmt.Sprintf("%v, %v", req.Method, req.URL)))

	now := time.Now()
	resp, err := t.next.RoundTrip(req)

	if err != nil {
		log.Error("Failed to send request", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("Request finished", slog.Duration(
		"duration", time.Since(now)), slog.Int("status", resp.StatusCode))
	return resp, nil
}

func (c *Client) GetTimetable(cityURL string) []time.Time {
	const op = "parser.Client.GetTimetable"
	log := c.log.With(slog.String("op", op))

	resp, err := c.Get(cityURL)
	if err != nil {
		log.Error("Failed to send request", slog.String("error", err.Error()))
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error("Failed to close response body", slog.String("error", err.Error()))
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error("Failed to load document", slog.String("error", err.Error()))
	}

	var timetable []time.Time

	columnCount := 0
	doc.Find("tr.active_day td").Each(func(i int, tr *goquery.Selection) {
		columnCount++
	})

	var requiredIndexes []int

	switch columnCount {
	case 7:
		requiredIndexes = []int{1, 3, 4, 5, 6}
	case 9:
		requiredIndexes = []int{2, 4, 5, 6, 8}
	}

	doc.Find("tr.active_day td").Each(func(i int, td *goquery.Selection) {
		if slices.Contains(requiredIndexes, i) {
			timeText := td.Find("b").Text()
			inTimeFormat, _ := time.Parse("15:04", timeText)
			now := time.Now()
			datetimeFormat := time.Date(now.Year(), now.Month(), now.Day(), inTimeFormat.Hour(),
				inTimeFormat.Minute(), 0, 0, now.Location())
			timetable = append(timetable, datetimeFormat)
		}
	})
	return timetable
}
