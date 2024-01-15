package parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"slices"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"

type roundTripper struct {
	next http.RoundTripper
}

type Client struct {
	*http.Client
}

func NewClient(timeout time.Duration) (*Client, error) {
	if timeout == 0 {
		return nil, errors.New("timeout must be greater than 0")
	}
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
			Transport: &roundTripper{
				next: http.DefaultTransport,
			},
		},
	}, nil
}

func (t *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	log.Printf("Making request: [%s] %s", req.Method, req.URL)
	now := time.Now()
	resp, err := t.next.RoundTrip(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return nil, err
	}
	log.Printf("Request successful. Status: %s. Elapsed time: %s", resp.Status, time.Since(now))
	return resp, nil
}

func (c *Client) GetTimetable(cityURL string) []time.Time {
	resp, err := c.Get(cityURL)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var timetable []time.Time

	requiredIndexes := []int{1, 3, 4, 5, 6}
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
