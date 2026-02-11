package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertTimeToNetherlandsTimezone(t time.Time) time.Time {
	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		return t
	}
	return t.In(loc)
}

func NullableDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func DatePtr(date pgtype.Date) *time.Time {
	if !date.Valid {
		return nil
	}
	t := date.Time
	return &t
}
