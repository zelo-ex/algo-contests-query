package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"
)

type LuoguContest struct {
	ID                 int    `json:"id"`
	StartTime          int64  `json:"startTime"`
	EndTime            int64  `json:"endTime"`
	Name               string `json:"name"`
	Method             int    `json:"method"`
	Visibility         int    `json:"visibility"`
	InvitationCodeType int    `json:"invitationCodeType"`
	Rated              int    `json:"rated"`
	Host               struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Ispremium bool   `json:"isPremium"`
	} `json:"host"`
	Squad        bool `json:"squad"`
	ProblemCount int  `json:"problemCount"`
}

type LuoguContestsList struct {
	Instance string `json:""`
	Templete string `json:""`
	Status   int    `json:""`
	Locale   string `json:""`
	Data     struct {
		Contests struct {
			PerPage int            `json:"perPage"`
			Count   int64          `json:"count"`
			Result  []LuoguContest `json:"result"`
		} `json:"contests"`
	} `json:"data"`
}

func GetLuoguContestsList() (ContestsList, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	req, err := http.NewRequest(http.MethodGet,
		"https://www.luogu.com.cn/contest/list", nil)
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

	regJSON, err := regexp.Compile(`(?s)<script id="lentille-context" type="application/json">(.*?)</script>`)
	if err != nil {
		return nil, err
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s := regJSON.FindStringSubmatch(string(html))

	var data LuoguContestsList
	err = json.NewDecoder(strings.NewReader(s[1])).Decode(&data)
	if err != nil {
		return nil, err
	}

	result := ContestsListConvert(
		data.Data.Contests.Result,
		func(t LuoguContest) Contest {
			c := Contest{
				Platform: "Luogu",
				Name:     t.Name,
				Start:    time.Unix(t.StartTime, 0),
				Duration: time.Duration(t.EndTime-t.StartTime) * time.Second,
			}
			return c
		},
	).Filter(func(c Contest) bool {
		return c.Start.Add(c.Duration).After(time.Now())
	}).Sort()

	return result, nil
}
