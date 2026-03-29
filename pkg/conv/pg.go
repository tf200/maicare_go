package conv

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TimeFromPgTimestamptz(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func StringFromPgTime(value pgtype.Time) string {
	if !value.Valid {
		return ""
	}

	totalSeconds := value.Microseconds / 1000000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
