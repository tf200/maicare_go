package util

import (
	"fmt"
	"strings"
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

func ParseYYYYMMDD(s string) (time.Time, bool, error) {
	if strings.TrimSpace(s) == "" {
		return time.Time{}, false, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid date %q (expected YYYY-MM-DD): %w", s, err)
	}
	return t, true, nil
}

func ParseRFC3339OrYYYYMMDD(s string) (time.Time, bool, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false, nil
	}
	if strings.Contains(s, "T") {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Time{}, false, fmt.Errorf("invalid timestamp %q (expected RFC3339): %w", s, err)
		}
		return t, true, nil
	}
	return ParseYYYYMMDD(s)
}
