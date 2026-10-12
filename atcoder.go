package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parseAtcoderDuration(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration: %q", s)
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}

func GetAtcoderContestsList() (ContestsList, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet,
		"https://atcoder.jp/contests/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	body := html.Find("#contest-table-upcoming .table tbody")
	contests := ContestsList{}
	body.Find("tr").Each(func(idxTr int, tr *goquery.Selection) {
		tds := tr.Find("td")

		start := tds.Eq(0).Find("time").Text()
		startTime, err := time.Parse("2006-01-02 15:04:05-0700", start)
		if err != nil {
			log.Println(err)
			return
		}

		title := tds.Eq(1).Find("a").Text()
		duration, err := parseAtcoderDuration(tds.Eq(2).Text())
		if err != nil {
			log.Println(err)
			return
		}

		contests = append(contests, Contest{
			Platform: "Atcoder",
			Name:     title,
			Start:    startTime.In(time.Local),
			Duration: duration,
		})
	})

	return contests, nil
}
