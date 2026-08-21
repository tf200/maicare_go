package service

import (
	"context"
	"testing"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

type medicalAuditRecorder struct {
	events []domain.AuditEvent
}

func (r *medicalAuditRecorder) Log(_ context.Context, event domain.AuditEvent) error {
	r.events = append(r.events, event)
	return nil
}

func (r *medicalAuditRecorder) LogBatch(_ context.Context, _ domain.AuditBatchEvent) error {
	return nil
}

func TestLogMedicalAuditExcludesClinicalPayload(t *testing.T) {
	recorder := &medicalAuditRecorder{}
	service := &ClientService{audit: recorder}
	clientID := uuid.New()
	subjectID := uuid.New()

	service.logMedicalAudit(
		context.Background(),
		"read",
		"client_diagnosis",
		subjectID,
		clientID,
		"CLIENT.DIAGNOSIS.VIEW",
		1,
	)

	if len(recorder.events) != 1 {
		t.Fatalf("audit event count = %d, want 1", len(recorder.events))
	}
	event := recorder.events[0]
	if event.EventType != "record_access" || event.Action != "read" || event.Result != "success" {
		t.Fatalf("unexpected medical audit classification: %+v", event)
	}
	if event.SubjectType != "client_diagnosis" || event.SubjectID != subjectID.String() {
		t.Fatalf("unexpected medical audit subject: %+v", event)
	}
	if event.ClientID == nil || *event.ClientID != clientID.String() {
		t.Fatalf("medical audit client ID = %v, want %s", event.ClientID, clientID)
	}
	if event.AccessRule == nil || *event.AccessRule != "CLIENT.DIAGNOSIS.VIEW" {
		t.Fatalf("medical audit access rule = %v", event.AccessRule)
	}
	if len(event.Details) != 1 || event.Details["count"] != 1 {
		t.Fatalf("medical audit details = %#v, want count only", event.Details)
	}
}
