package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type CFContest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	StartTimeSeconds int64  `json:"startTimeSeconds"`
	DurationSeconds  int64  `json:"durationSeconds"`
}

type CFContestsList struct {
	Status string      `json:"status"`
	Result []CFContest `json:"result"`
}

func GetCodeForcesContestsList() (ContestsList, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet,
		"https://codeforces.com/api/contest.list?gym=false", nil)
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

	var data CFContestsList
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	result := ContestsListConvert(data.Result, func(t CFContest) Contest {
		c := Contest{
			Platform: "CodeForces",
			Name:     t.Name,
			Start:    time.Unix(t.StartTimeSeconds, 0),
			Duration: time.Duration(t.DurationSeconds) * time.Second,
		}
		return c
	}).Filter(func(c Contest) bool {
		return (c.Start.Add(c.Duration).After(time.Now()))
	}).Sort()

	return result, nil
}
