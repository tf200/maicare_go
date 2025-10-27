package audit

import (
	"context"
	db "maicare_go/db/sqlc"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuditService struct {
	Store *db.Store
}

func NewAuditService(store *db.Store) *AuditService {
	return &AuditService{
		Store: store,
	}
}

type AuditRecord struct {
	EventID      uuid.UUID
	EventType    string
	OccuredAt    time.Time
	ActorRole    string
	ActorID      int64
	SubjectType  string
	SubjectID    int64
	Action       string
	Result       string
	AccessReason string
	Ip           *netip.Addr
	SelfHash     string
	PreviousHash string
}

func (s *AuditService) CreateAuditRecord(ctx context.Context, record *AuditRecord) error {
	return s.Store.CreateAuditRecord(ctx, db.CreateAuditRecordParams{
		EventID:      record.EventID,
		EventType:    record.EventType,
		OccuredAt:    pgtype.Timestamptz{Time: record.OccuredAt, Valid: true},
		ActorRole:    record.ActorRole,
		ActorID:      record.ActorID,
		SubjectType:  record.SubjectType,
		SubjectID:    record.SubjectID,
		Action:       record.Action,
		Result:       record.Result,
		AccessReason: record.AccessReason,
		Ip:           record.Ip,
		HashPrev:     record.PreviousHash,
		HashSelf:     record.SelfHash,
	})
}

func (s *AuditService) GetLatestAuditRecordHash(ctx context.Context, subjectID int64) (string, error) {
	return s.Store.GetLatestAuditHashBySubject(ctx, subjectID)

	