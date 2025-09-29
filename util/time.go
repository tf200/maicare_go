package util

import (
	"time"
)

func ConvertTimeToNetherlandsTimezone(t time.Time) time.Time {
	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		return t
	}
	return t.In(loc)
}
