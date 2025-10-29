package audit

import (
	"context"
	"database/sql"
	"errors"
	db "maicare_go/db/sqlc"

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

// service/audit/service.go
func (s *AuditService) CreateAuditRecord(ctx context.Context, record *AuditRecord) error {
	lastHash, err := s.Store.GetLatestAuditHash(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	selfHash := record.calculateAuditHash()
	record.SelfHash = selfHash
	record.PreviousHash = lastHash

	err = record.validate()
	if err != nil {
		return err
	}

	s.Store.CreateAuditRecord(ctx, db.CreateAuditRecordParams{
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

	return nil
}
