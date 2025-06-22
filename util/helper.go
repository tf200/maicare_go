package util

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int64) *int64 {
	return &i
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func Float64Ptr(f float64) *float64 {
	return &f

}

// Helper function to convert string to pgtype.Time
func StringToPgTime(timeStr string) (pgtype.Time, error) {
	// Parse the time string (expects format like "10:00:00" or "10:00")
	var t time.Time
	var err error

	// Try parsing with seconds first
	t, err = time.Parse("15:04:05", timeStr)
	if err != nil {
		// Try parsing without seconds
		t, err = time.Parse("15:04", timeStr)
		if err != nil {
			return pgtype.Time{}, err
		}
	}

	// Convert to microseconds since midnight
	microseconds := int64(t.Hour())*3600000000 +
		int64(t.Minute())*60000000 +
		int64(t.Second())*1000000

	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}, nil
}

// In your util package
func PgTimeToString(pgTime pgtype.Time) string {
	if !pgTime.Valid {
		return ""
	}

	// Convert microseconds to hours, minutes, seconds
	totalSeconds := pgTime.Microseconds / 1000000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func MicrosecondsToTimeComponents(microseconds int64) (hour, min, sec, nanosec int) {
	const (
		microsecondsPerHour   = 3600000000
		microsecondsPerMinute = 60000000
		microsecondsPerSecond = 1000000
	)

	remaining := microseconds

	hour = int(remaining / microsecondsPerHour)
	remaining = remaining % microsecondsPerHour

	min = int(remaining / microsecondsPerMinute)
	remaining = remaining % microsecondsPerMinute

	sec = int(remaining / microsecondsPerSecond)
	remaining = remaining % microsecondsPerSecond

	nanosec = int(remaining * 1000) // Convert remaining microseconds to nanoseconds

	return hour, min, sec, nanosec
}

func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func DerefInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func DerefUUID(u *uuid.UUID) string {
	if u == nil {
		return ""
	}
	return u.String()
}
