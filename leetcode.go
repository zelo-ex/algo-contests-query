package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type LCContest struct {
	Title     string `json:"title"`
	TitleCn   string `json:"titleCn"`
	StartTime int64  `json:"startTime"`
	Duration  int    `json:"duration"`
}

type LCContestsList struct {
	Data struct {
		Contests []LCContest `json:"contestV2UpcomingContests"`
	} `json:"data"`
}

func GetLeetcodeContestsList() (ContestsList, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}
	_, err = client.Get("https://leetcode.cn/contest/")
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse("https://leetcode.cn")
	var csrfToken string
	for _, c := range client.Jar.Cookies(u) {
		if c.Name == "csrftoken" {
			csrfToken = c.Value
		}
	}

	payload := struct {
		Query         string         `json:"query"`
		Variables     map[string]any `json:"variables,omitempty"`
		OperationName string         `json:"operationName,omitempty"`
	}{
		Query: `query contestV2UpcomingContests {
	contestV2UpcomingContests {
		title
		titleCn
		startTime
		duration
	}
}`,
		OperationName: "contestV2UpcomingContests",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://leetcode.cn/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrfToken})
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrftoken", csrfToken)
	req.Header.Set("x-operation-name", "contestV2UpcomingContests")
	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	req.Header.Set("Origin", "https://leetcode.cn")
	req.Header.Set("Referer", "https://leetcode.cn/contest/")

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

	var data LCContestsList
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	result := ContestsListConvert(
		data.Data.Contests,
		func(t LCContest) Contest {
			c := Contest{
				Platform: "Leetcode",
				Name:     t.TitleCn,
				Start:    time.Unix(t.StartTime, 0),
				Duration: time.Duration(t.Duration) * time.Second,
			}
			return c
		},
	).Filter(func(c Contest) bool {
		return c.Start.Add(c.Duration).After(time.Now())
	}).Sort()

	return result, nil
}
