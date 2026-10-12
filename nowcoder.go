package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type NCContest struct {
	ContestStartTime int64  `json:"contestStartTime"`
	ContestDuration  int64  `json:"contestDuration"`
	ContestName      string `json:"contestName"`
}

func GetNowcoderContestsList() (ContestsList, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet,
		"https://ac.nowcoder.com/acm/contest/vip-index", nil)
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []NCContest
	body := doc.Find("div.platform-mod.js-current>.js-item")
	body.Each(func(i int, s *goquery.Selection) {
		itemAttr, isExist := s.Attr("data-json")
		if !isExist {
			return
		}
		itemAttr = html.UnescapeString(itemAttr)
		var curr NCContest
		err := json.NewDecoder(strings.NewReader(itemAttr)).Decode(&curr)
		if err != nil {
			log.Println(err)
			return
		}
		data = append(data, curr)
	})

	result := ContestsListConvert(data, func(t NCContest) Contest {
		duration := time.Duration(t.ContestDuration) * time.Millisecond
		c := Contest{
			Platform: "Nowcoder",
			Name:     t.ContestName,
			Start:    time.UnixMilli(t.ContestStartTime),
			Duration: duration,
		}
		return c
	}).Sort()

	return result, nil
}
