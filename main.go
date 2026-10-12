package main

import (
	"fmt"
	"log"
)

func main() {
	type Result struct {
		contests ContestsList
		err      error
	}

	fetchers := []func() (ContestsList, error){
		GetCodeForcesContestsList,
		GetLeetcodeContestsList,
		GetLuoguContestsList,
		GetAtcoderContestsList,
		GetNowcoderContestsList,
	}

	ch := make(chan Result, len(fetchers))
	for _, f := range fetchers {
		go func() {
			c, e := f()
			ch <- Result{c, e}
		}()
	}

	var contests ContestsList
	for range len(fetchers) {
		r := <-ch
		if r.err != nil {
			log.Println(r.err)
		}
		contests = append(contests, r.contests...)
	}

	contests = contests.Sort()
	for _, t := range contests {
		fmt.Printf("platform : %s\nname     : %s\n"+
			"time     : %s\nduration : %s\n\n",
			t.Platform, t.Name,
			t.Start.Format("2006-01-02 15:04:05"),
			DurationFormat(t.Duration),
		)
	}
}
