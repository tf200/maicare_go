package invoice

import (
	"fmt"
	"math"
	"time"
)

const (
	DefaultBillingTimezone = "Europe/Amsterdam"
	DefaultBillingCycle    = "iso_4_week"
)

func centsFromAmount(amount float64) int64 {
	// Half away from zero; consistent for invoices.
	return int64(math.Round(amount * 100))
}

func amountFromCents(cents int64) float64 {
	return float64(cents) / 100
}

func roundHalfUpDiv(n, d int64) int64 {
	if d == 0 {
		panic("division by zero")
	}
	if n == 0 {
		return 0
	}
	sign := int64(1)
	if n < 0 {
		sign = -1
		n = -n
	}
	q := n / d
	r := n % d
	if 2*r >= d {
		q++
	}
	return sign * q
}

func clampTime(t, min, max time.Time) time.Time {
	if t.Before(min) {
		return min
	}
	if t.After(max) {
		return max
	}
	return t
}

func countCalendarDays(start, end time.Time, tz string) (int64, error) {
	if end.Before(start) || end.Equal(start) {
		return 0, nil
	}
	if tz == "" {
		tz = DefaultBillingTimezone
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return 0, fmt.Errorf("invalid billing timezone %q: %w", tz, err)
	}

	sy, sm, sd := start.In(loc).Date()
	ey, em, ed := end.In(loc).Date()

	// Use UTC midnights for the local dates to avoid DST hour-length anomalies.
	sUTC := time.Date(sy, sm, sd, 0, 0, 0, 0, time.UTC)
	eUTC := time.Date(ey, em, ed, 0, 0, 0, 0, time.UTC)
	if eUTC.Before(sUTC) {
		return 0, nil
	}
	return int64(eUTC.Sub(sUTC).Hours() / 24), nil
}

type lineAmounts struct {
	netCents   int64
	vatCents   int64
	grossCents int64
}

func computeVat(netCents int64, vatRatePct int64) lineAmounts {
	vatCents := roundHalfUpDiv(netCents*vatRatePct, 100)
	return lineAmounts{
		netCents:   netCents,
		vatCents:   vatCents,
		grossCents: netCents + vatCents,
	}
}

func roundTo(amount float64, decimals int) float64 {
	pow := math.Pow10(decimals)
	return math.Round(amount*pow) / pow
}
