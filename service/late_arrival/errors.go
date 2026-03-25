package late_arrival

import "errors"

var (
	ErrLateArrivalInvalidRequest = errors.New("invalid late arrival request")
	ErrLateArrivalConflict       = errors.New("late arrival conflict")
)
