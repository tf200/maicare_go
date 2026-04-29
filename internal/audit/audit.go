// Package audit implements NEN 7513-compliant audit logging.
//
// It reads actor/request metadata from context and manages the hash chain for
// tamper-evident audit records. All hash-chain logic is contained here.
package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/netip"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/internal/middleware"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// Querier is the subset of db.Queries that the audit package needs.
// This makes testing easier and reduces coupling.
type Querier interface {
	CreateAuditRecord(ctx context.Context, arg db.CreateAuditRecordParams) error
	BulkCreateAuditRecords(ctx context.Context, arg []db.BulkCreateAuditRecordsParams) (int64, error)
	GetLatestAuditHash(ctx context.Context) (string, error)
	LockAuditHashChain(ctx context.Context) error
}

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn db.TxFn) error
}

type auditLogger struct {
	store  Store
	logger domain.Logger
}

// New creates a new AuditLogger.
func New(store Store, logger domain.Logger) domain.AuditLogger {
	return &auditLogger{store: store, logger: logger}
}

// Log writes a single audit record.
// Actor/request metadata is extracted from ctx automatically.
func (l *auditLogger) Log(ctx context.Context, event domain.AuditEvent) error {
	err := l.store.ExecTx(ctx, func(q *db.Queries) error {
		if err := q.LockAuditHashChain(ctx); err != nil {
			return err
		}
		prevHash := fetchLatestHash(ctx, q)
		row, err := l.buildRow(ctx, event, nil, prevHash)
		if err != nil {
			return err
		}
		return q.CreateAuditRecord(ctx, row)
	})
	if err != nil {
		l.logger.LogError(ctx, "AuditLogger.Log", "failed to write audit record", err,
			zap.String("event_type", event.EventType),
			zap.String("action", event.Action),
			zap.String("subject_type", event.SubjectType),
			zap.String("subject_id", event.SubjectID),
		)
		return err
	}
	return nil
}

// LogBatch writes multiple audit records in one COPY batch.
// All rows share common metadata and event_group_id.
// Hashes are computed sequentially in application memory before the DB write.
func (l *auditLogger) LogBatch(ctx context.Context, event domain.AuditBatchEvent) error {
	if len(event.SubjectIDs) == 0 {
		return nil
	}

	groupID, err := uuid.Parse(event.EventGroupID)
	if err != nil {
		return fmt.Errorf("invalid audit event_group_id %q: %w", event.EventGroupID, err)
	}

	err = l.store.ExecTx(ctx, func(q *db.Queries) error {
		if err := q.LockAuditHashChain(ctx); err != nil {
			return err
		}

		now := time.Now()
		meta := l.extractMetadata(ctx)
		hashPrev := fetchLatestHash(ctx, q)

		rows := make([]db.BulkCreateAuditRecordsParams, 0, len(event.SubjectIDs))
		for _, sid := range event.SubjectIDs {
			subjID := sid

			// Default client_id to subject_id when subject is a client record
			var clientID *uuid.UUID
			if event.ClientID != nil {
				cid, err := uuid.Parse(*event.ClientID)
				if err != nil {
					return fmt.Errorf("invalid audit client_id %q: %w", *event.ClientID, err)
				}
				clientID = &cid
			} else if event.SubjectType == "client" {
				cid, err := uuid.Parse(subjID)
				if err == nil {
					clientID = &cid
				}
			}

			detailsJSON := marshalDetails(event.Details)

			row := db.BulkCreateAuditRecordsParams{
				EventID:         uuid.New(),
				EventGroupID:    &groupID,
				OccurredAt:      pgtype.Timestamptz{Time: now, Valid: true},
				EventType:       event.EventType,
				Action:          event.Action,
				Result:          event.Result,
				ActorUserID:     meta.actorUserID,
				ActorEmployeeID: meta.actorEmployeeID,
				ActorRoles:      meta.actorRoles,
				SubjectType:     event.SubjectType,
				SubjectID:       &subjID,
				ClientID:        clientID,
				AccessRule:      meta.accessRule(event.AccessRule),
				AccessReason:    event.AccessReason,
				SessionID:       meta.sessionID,
				RequestID:       meta.requestID,
				Ip:              meta.ip,
				UserAgent:       meta.userAgent,
				Route:           meta.route,
				Method:          meta.method,
				Details:         detailsJSON,
			}

			row.HashPrev = hashPrev
			row.HashSelf = computeHashSelf(row, hashPrev)
			hashPrev = row.HashSelf

			rows = append(rows, row)
		}

		_, err := q.BulkCreateAuditRecords(ctx, rows)
		return err
	})
	if err != nil {
		l.logger.LogError(ctx, "AuditLogger.LogBatch", "failed to write batch audit records", err,
			zap.String("event_type", event.EventType),
			zap.String("action", event.Action),
			zap.Int("count", len(event.SubjectIDs)),
		)
		return err
	}
	return nil
}

// buildRow constructs a single CreateAuditRecordParams from a domain event.
func (l *auditLogger) buildRow(ctx context.Context, event domain.AuditEvent, groupID *uuid.UUID, prevHash string) (db.CreateAuditRecordParams, error) {
	now := time.Now()
	meta := l.extractMetadata(ctx)

	subjID := event.SubjectID

	// Default client_id to subject_id when subject is a client record
	var clientID *uuid.UUID
	if event.ClientID != nil {
		cid, err := uuid.Parse(*event.ClientID)
		if err != nil {
			return db.CreateAuditRecordParams{}, fmt.Errorf("invalid audit client_id %q: %w", *event.ClientID, err)
		}
		clientID = &cid
	} else if event.SubjectType == "client" {
		cid, err := uuid.Parse(subjID)
		if err == nil {
			clientID = &cid
		}
	}

	detailsJSON := marshalDetails(event.Details)

	row := db.CreateAuditRecordParams{
		EventID:         uuid.New(),
		EventGroupID:    groupID,
		OccurredAt:      pgtype.Timestamptz{Time: now, Valid: true},
		EventType:       event.EventType,
		Action:          event.Action,
		Result:          event.Result,
		ActorUserID:     meta.actorUserID,
		ActorEmployeeID: meta.actorEmployeeID,
		ActorRoles:      meta.actorRoles,
		SubjectType:     event.SubjectType,
		SubjectID:       &subjID,
		ClientID:        clientID,
		AccessRule:      meta.accessRule(event.AccessRule),
		AccessReason:    event.AccessReason,
		SessionID:       meta.sessionID,
		RequestID:       meta.requestID,
		Ip:              meta.ip,
		UserAgent:       meta.userAgent,
		Route:           meta.route,
		Method:          meta.method,
		Details:         detailsJSON,
	}

	row.HashPrev = prevHash
	row.HashSelf = computeHashSelfFromCreateParams(row, prevHash)

	return row, nil
}

// metadata holds all request/actor information extracted from context.
type metadata struct {
	actorUserID     *uuid.UUID
	actorEmployeeID *uuid.UUID
	actorRoles      []string
	sessionID       *uuid.UUID
	requestID       *string
	ip              *netip.Addr
	userAgent       *string
	route           *string
	method          *string
}

func (m metadata) accessRule(rule *string) *string {
	if rule != nil {
		return rule
	}
	return nil
}

// extractMetadata reads all relevant fields from the Go context.
// All fields are optional — the audit table allows NULLs for actor data
// (e.g. unauthenticated login attempts).
func (l *auditLogger) extractMetadata(ctx context.Context) metadata {
	var m metadata

	m.requestID = strPtr(middleware.RequestIDFromContextSafe(ctx))
	m.ip = middleware.ClientIPFromContextSafe(ctx)
	m.userAgent = strPtr(middleware.UserAgentFromContextSafe(ctx))
	m.route = strPtr(middleware.RouteFromContextSafe(ctx))
	m.method = strPtr(middleware.MethodFromContextSafe(ctx))

	// Actor info from auth payload
	if payload, ok := middleware.AuthPayloadFromContext(ctx); ok && payload != nil {
		uid := payload.UserID
		if uid != uuid.Nil {
			m.actorUserID = &uid
		}
		eid := payload.EmployeeID
		if eid != uuid.Nil {
			m.actorEmployeeID = &eid
		}
		sid := payload.SessionID
		if sid != uuid.Nil {
			m.sessionID = &sid
		}
	}

	// Actor roles from RBAC middleware
	if roles, ok := middleware.ActorRolesFromContext(ctx); ok {
		m.actorRoles = roles
	}

	return m
}

// fetchLatestHash retrieves the most recent hash from the audit table.
// If the table is empty, returns an empty string (genesis block).
func fetchLatestHash(ctx context.Context, q Querier) string {
	hash, err := q.GetLatestAuditHash(ctx)
	if err != nil {
		// Table empty or not yet seeded — start the chain
		return ""
	}
	return hash
}

// computeHashSelfFromCreateParams computes the SHA-256 hash for a single
// CreateAuditRecordParams row, including hash_prev.
func computeHashSelfFromCreateParams(row db.CreateAuditRecordParams, hashPrev string) string {
	return computeHash(hashPrev, fieldsFromCreateParams(row))
}

// computeHashSelf computes the SHA-256 hash for a BulkCreateAuditRecordsParams row.
func computeHashSelf(row db.BulkCreateAuditRecordsParams, hashPrev string) string {
	return computeHash(hashPrev, fieldsFromBulkParams(row))
}

// computeHash builds the SHA-256 hash from a list of field values.
// The hash chain is: SHA256(hash_prev || '|' || field1 || '|' || field2 || ...)
func computeHash(hashPrev string, fields []string) string {
	h := sha256.New()
	h.Write([]byte(hashPrev))
	for _, f := range fields {
		h.Write([]byte("|"))
		h.Write([]byte(f))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func fieldsFromCreateParams(row db.CreateAuditRecordParams) []string {
	return []string{
		row.EventID.String(),
		uuidOrNil(row.EventGroupID),
		row.OccurredAt.Time.Format(time.RFC3339Nano),
		row.EventType,
		row.Action,
		row.Result,
		uuidOrNil(row.ActorUserID),
		uuidOrNil(row.ActorEmployeeID),
		stringsJoin(row.ActorRoles),
		row.SubjectType,
		strOrEmpty(row.SubjectID),
		uuidOrNil(row.ClientID),
		strOrEmpty(row.AccessRule),
		strOrEmpty(row.AccessReason),
		uuidOrNil(row.SessionID),
		strOrEmpty(row.RequestID),
		ipOrEmpty(row.Ip),
		strOrEmpty(row.UserAgent),
		strOrEmpty(row.Route),
		strOrEmpty(row.Method),
		stringOrEmpty(row.Details),
	}
}

func fieldsFromBulkParams(row db.BulkCreateAuditRecordsParams) []string {
	return []string{
		row.EventID.String(),
		uuidOrNil(row.EventGroupID),
		row.OccurredAt.Time.Format(time.RFC3339Nano),
		row.EventType,
		row.Action,
		row.Result,
		uuidOrNil(row.ActorUserID),
		uuidOrNil(row.ActorEmployeeID),
		stringsJoin(row.ActorRoles),
		row.SubjectType,
		strOrEmpty(row.SubjectID),
		uuidOrNil(row.ClientID),
		strOrEmpty(row.AccessRule),
		strOrEmpty(row.AccessReason),
		uuidOrNil(row.SessionID),
		strOrEmpty(row.RequestID),
		ipOrEmpty(row.Ip),
		strOrEmpty(row.UserAgent),
		strOrEmpty(row.Route),
		strOrEmpty(row.Method),
		stringOrEmpty(row.Details),
	}
}

// --- helpers ---

func marshalDetails(details map[string]any) []byte {
	if len(details) == 0 {
		return nil
	}
	b, err := json.Marshal(details)
	if err != nil {
		return nil
	}
	return b
}

func uuidOrNil(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringOrEmpty(b []byte) string {
	if b == nil {
		return ""
	}
	return string(b)
}

func ipOrEmpty(addr *netip.Addr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}

func stringsJoin(s []string) string {
	return strings.Join(s, "|")
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Ensure auditLogger implements domain.AuditLogger.
var _ domain.AuditLogger = (*auditLogger)(nil)

// init registers the hash algorithm used. Useful for future audit verification.
const HashAlgorithm = "SHA-256"
