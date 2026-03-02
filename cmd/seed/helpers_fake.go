package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

func pgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func pgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func fakePhone() string {
	return fmt.Sprintf("06%08d", gofakeit.Number(0, 99999999))
}

func fakeBSN(index int) string {
	return fmt.Sprintf("%09d", (index%900000000)+100000000)
}

func fakePostalCodeNL() string {
	return fmt.Sprintf("%04d%s", gofakeit.Number(1000, 9999), strings.ToUpper(gofakeit.LetterN(2)))
}

func randomDate(startYear, endYear int) time.Time {
	start := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, 12, 31, 0, 0, 0, 0, time.UTC)
	return gofakeit.DateRange(start, end)
}

func randomRecentDate(maxDaysAgo int) time.Time {
	now := time.Now().UTC()
	start := now.AddDate(0, 0, -maxDaysAgo)
	return gofakeit.DateRange(start, now)
}

func sanitizeDomain(s string) string {
	clean := strings.ToLower(strings.TrimSpace(s))
	clean = strings.ReplaceAll(clean, "&", "and")
	clean = strings.ReplaceAll(clean, " ", "")
	if clean == "" {
		return "company"
	}
	return clean
}

func chance(probability float64) bool {
	return gofakeit.Float64() < probability
}

func nullableString(value string, nilProbability float64) *string {
	if chance(nilProbability) {
		return nil
	}
	return &value
}
