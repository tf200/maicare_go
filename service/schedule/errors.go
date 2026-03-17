package schedule

import "errors"

var ErrWeekNotEmpty = errors.New("week is not empty")

var (
	ErrScheduleNotFound                = errors.New("schedule not found")
	ErrShiftSwapNotFound               = errors.New("shift swap request not found")
	ErrShiftSwapInvalidRequest         = errors.New("invalid shift swap request")
	ErrShiftSwapStateInvalid           = errors.New("shift swap request is not in a valid state")
	ErrShiftSwapExpired                = errors.New("shift swap request has expired")
	ErrShiftSwapConflict               = errors.New("swap would create schedule overlap conflict")
	ErrShiftSwapScheduleOwnership      = errors.New("schedule ownership is invalid for this swap")
	ErrShiftSwapDuplicateActiveRequest = errors.New("one of the schedules is already in an active swap request")
)
