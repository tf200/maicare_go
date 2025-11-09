package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AuditRecord struct {
	EventID      uuid.UUID   `json:"event_id"`
	EventType    string      `json:"event_type"`
	OccuredAt    time.Time   `json:"occured_at"`
	ActorRole    []string    `json:"actor_role"`
	ActorID      uuid.UUID   `json:"actor_id"`
	SubjectType  string      `json:"subject_type"`
	SubjectID    uuid.UUID   `json:"subject_id"`
	Action       string      `json:"action"`
	Result       string      `json:"result"`
	AccessReason string      `json:"access_reason"`
	Ip           *netip.Addr `json:"ip"`
	SelfHash     string      `json:"self_hash"`
	PreviousHash string      `json:"previous_hash"`
}

// calculateAuditHash generates a hash for the audit record for integrity verification
func (r *AuditRecord) calculateAuditHash() string {
	// Create a deterministic string from all audit record fields
	ipStr := ""
	if r.Ip != nil {
		ipStr = r.Ip.String()
	}

	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		r.EventID.String(),
		r.EventType,
		r.OccuredAt.Format(time.RFC3339Nano),
		r.ActorID.String(),
		r.ActorRole,
		r.SubjectType,
		r.SubjectID.String(),
		r.Action,
		r.Result,
		r.AccessReason,
		ipStr,
	)

	// Calculate SHA256 hash
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *AuditRecord) validate() error {
	var errs []string

	if s.EventID == uuid.Nil {
		errs = append(errs, "EventID is required")
	}

	if s.EventType == "" {
		errs = append(errs, "EventType is required")
	}

	if s.OccuredAt.IsZero() {
		errs = append(errs, "OccuredAt is required")
	}

	if len(s.ActorRole) == 0 {
		errs = append(errs, "ActorRole is required")
	}

	if s.ActorID == uuid.Nil {
		errs = append(errs, "ActorID is required")
	}

	if s.SubjectType == "" {
		errs = append(errs, "SubjectType is required")
	}

	if s.SubjectID == uuid.Nil {
		errs = append(errs, "SubjectID is required")
	}

	if s.Action == "" {
		errs = append(errs, "Action is required")
	}

	if s.Result == "" {
		errs = append(errs, "Result is required")
	}

	if s.AccessReason == "" {
		errs = append(errs, "AccessReason is required")
	}

	if s.Ip == nil {
		errs = append(errs, "Ip is required")
	}

	if s.SelfHash == "" {
		errs = append(errs, "SelfHash is required")
	}

	if s.PreviousHash == "" {
		errs = append(errs, "PreviousHash is required")
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}

	return nil
}
