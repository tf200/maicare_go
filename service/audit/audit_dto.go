package audit

import (
	"maicare_go/pagination"
	"time"

	"github.com/google/uuid"
)

type ListAuditRecordsRequest struct {
	SubjectID *uuid.UUID `form:"subject_id"`
	ActorID   *uuid.UUID `form:"actor_id"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"`
	*pagination.Request
}


