package domain

import "context"

// AuditEvent carries semantic business information for a single audit record.
// The internal/audit package maps this to the canonical DB schema, enriches it
// with actor/request metadata from context, and computes the hash chain.
type AuditEvent struct {
	EventType   string         // record_access, authentication, export, etc.
	Action      string         // read, create, update, delete, list, export, login, logout
	Result      string         // success, failure, blocked
	SubjectType string         // client, registration_form, intake_form, user, session
	SubjectID   string         // the record ID being accessed
	ClientID    *string        // owning client ID, nil if not applicable
	AccessRule  *string        // permission name like "CLIENT.VIEW"
	AccessReason *string       // reason for access, nil if not collected
	Details     map[string]any // arbitrary structured data (filter params, error details, etc.)
}

// AuditBatchEvent describes a batch of audit rows sharing an event_group_id.
// One audit row is inserted per SubjectID, all sharing the same metadata.
type AuditBatchEvent struct {
	EventGroupID string
	EventType    string
	Action       string
	Result       string
	SubjectType  string
	ClientID     *string
	AccessRule   *string
	AccessReason *string
	SubjectIDs   []string
	Details      map[string]any
}

// AuditLogger is the domain contract for writing NEN 7513 audit records.
// Implementations read actor/request metadata from context and handle
// hash-chain computation internally.
type AuditLogger interface {
	// Log writes a single audit record.
	// Metadata (actor, IP, request_id, etc.) is extracted from ctx.
	Log(ctx context.Context, event AuditEvent) error

	// LogBatch writes multiple audit records in one batch operation.
	// All rows share the same event_group_id and metadata.
	// Metadata is extracted from ctx.
	LogBatch(ctx context.Context, event AuditBatchEvent) error
}
