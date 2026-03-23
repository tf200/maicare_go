package employees

import (
	"time"

	"github.com/google/uuid"
)

const (
	TimelineItemTypeShift = "shift"
	TimelineItemTypeEvent = "event"
)

type GetMyScheduleTimelineRequest struct {
	StartDate string `form:"start_date" binding:"required" example:"2026-02-01"`
	EndDate   string `form:"end_date" binding:"required" example:"2026-02-29"`
}

type GetMyScheduleTimelineDayResponse struct {
	Date  string                      `json:"date"`
	Items []GetMyScheduleTimelineItem `json:"items"`
}

type GetMyScheduleTimelineItem struct {
	ItemType  string                      `json:"item_type"`
	StartTime time.Time                   `json:"start_time"`
	EndTime   time.Time                   `json:"end_time"`
	Shift     *GetMyScheduleShiftItemInfo `json:"shift,omitempty"`
	Event     *GetMyScheduleEventItemInfo `json:"event,omitempty"`
}

type GetMyScheduleShiftItemInfo struct {
	ScheduleID   uuid.UUID `json:"schedule_id"`
	LocationID   uuid.UUID `json:"location_id"`
	LocationName string    `json:"location_name"`
}

type GetMyScheduleEventItemInfo struct {
	EventID            uuid.UUID  `json:"event_id"`
	MasterEventID      *uuid.UUID `json:"master_event_id,omitempty"`
	Title              string     `json:"title"`
	Description        *string    `json:"description,omitempty"`
	Location           *string    `json:"location,omitempty"`
	Color              *string    `json:"color,omitempty"`
	WorkApprovalStatus string     `json:"work_approval_status"`
	RecurrenceID       *time.Time `json:"recurrence_id,omitempty"`
}
