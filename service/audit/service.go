package audit

import (
	"context"
	"database/sql"
	"errors"
	db "maicare_go/db/sqlc"

	"github.com/gin-gonic/gin"
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
	if errors.Is(err, sql.ErrNoRows) {
		lastHash = ""
	} else if err != nil {
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

func (s *AuditService) ListAuditRecords(ctx *gin.Context, req ListAuditRecordsRequest) ([]AuditRecord, error) {
	params := req.GetParams()

	args := db.ListAuditRecordsParams{
		Limit:     params.Limit,
		Offset:    params.Offset,
		SubjectID: req.SubjectID,
		ActorID:   req.ActorID,
	}
	if req.StartTime != nil {
		args.StartTime = pgtype.Timestamptz{Time: *req.StartTime, Valid: true}
	}
	if req.EndTime != nil {
		args.EndTime = pgtype.Timestamptz{Time: *req.EndTime, Valid: true}
	}

	dbRecords, err := s.Store.ListAuditRecords(ctx, args)
	if err != nil {
		return nil, err
	}

	var records []AuditRecord
	for _, dbRec := range dbRecords {

		record := AuditRecord{
			EventID:      dbRec.EventID,
			EventType:    dbRec.EventType,
			OccuredAt:    dbRec.OccuredAt.Time,
			ActorRole:    dbRec.ActorRole,
			ActorID:      dbRec.ActorID,
			SubjectType:  dbRec.SubjectType,
			SubjectID:    dbRec.SubjectID,
			Action:       dbRec.Action,
			Result:       dbRec.Result,
			AccessReason: dbRec.AccessReason,
			Ip:           dbRec.Ip,
			SelfHash:     dbRec.HashSelf,
			PreviousHash: dbRec.HashPrev,
		}
		records = append(records, record)
	}

	return records, nil

}
