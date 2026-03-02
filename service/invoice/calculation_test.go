package invoice

import (
	"testing"
	"time"
)

func TestRoundHalfUpDiv(t *testing.T) {
	tests := []struct {
		n, d int64
		want int64
	}{
		{0, 7, 0},
		{1, 2, 1},
		{2, 2, 1},
		{3, 2, 2},
		{10, 4, 3},  // 2.5 -> 3
		{-10, 4, -3}, // -2.5 -> -3
	}
	for _, tc := range tests {
		if got := roundHalfUpDiv(tc.n, tc.d); got != tc.want {
			t.Fatalf("roundHalfUpDiv(%d,%d)=%d want %d", tc.n, tc.d, got, tc.want)
		}
	}
}

func TestCountCalendarDays_DSTSafe(t *testing.T) {
	// Use a timezone with DST to ensure we count calendar days, not hours.
	start := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC)
	days, err := countCalendarDays(start, end, "Europe/Amsterdam")
	if err != nil {
		t.Fatalf("countCalendarDays error: %v", err)
	}
	if days != 2 {
		t.Fatalf("countCalendarDays got %d want 2", days)
	}
}

func TestComputeVat(t *testing.T) {
	amts := computeVat(1000, 20)
	if amts.netCents != 1000 || amts.vatCents != 200 || amts.grossCents != 1200 {
		t.Fatalf("computeVat unexpected: %+v", amts)
	}
}

