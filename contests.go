package main

import (
	"fmt"
	"slices"
	"time"
)

const DEFAULT_USER_AGENT string = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36"

type Contest struct {
	Platform string
	Name     string
	Start    time.Time
	Duration time.Duration
}

type ContestsList []Contest

func (items ContestsList) Filter(keep func(Contest) bool) ContestsList {
	out := make([]Contest, 0, len(items))
	for _, it := range items {
		if keep(it) {
			out = append(out, it)
		}
	}
	return out
}

func ContestsListConvert[T any](t []T, conv func(T) Contest) ContestsList {
	out := make([]Contest, 0, len(t))
	for _, it := range t {
		out = append(out, conv(it))
	}
	return out
}

func (items ContestsList) Sort() ContestsList {
	out := slices.Clone(items)
	slices.SortFunc(out, func(a, b Contest) int {
		return a.Start.Compare(b.Start)
	})
	return out
}

func DurationFormat(d time.Duration) string {
	days, hours := int(d.Hours())/24, int(d.Hours())%24
	minutes, seconds := int(d.Minutes())%60, int(d.Seconds())%60
	return fmt.Sprintf("%02d:%02d:%02d:%02d", days, hours, minutes, seconds)
}
