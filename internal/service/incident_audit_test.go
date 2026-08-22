package service

import (
	"context"
	"testing"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

type incidentAuditRecorder struct {
	events []domain.AuditEvent
}

func (r *incidentAuditRecorder) Log(_ context.Context, event domain.AuditEvent) error {
	r.events = append(r.events, event)
	return nil
}

func (r *incidentAuditRecorder) LogBatch(_ context.Context, _ domain.AuditBatchEvent) error {
	return nil
}

func TestLogIncidentAuditExcludesIncidentPayload(t *testing.T) {
	recorder := &incidentAuditRecorder{}
	service := &incidentService{audit: recorder}
	clientID := uuid.New()
	incidentID := uuid.New()

	service.logIncidentAudit(context.Background(), "export", incidentID.String(), &clientID, "CLIENT.INCIDENT.VIEW", -1)

	if len(recorder.events) != 1 {
		t.Fatalf("audit event count = %d, want 1", len(recorder.events))
	}
	event := recorder.events[0]
	if event.EventType != "export" || event.Action != "export" || event.Result != "success" {
		t.Fatalf("unexpected incident audit classification: %+v", event)
	}
	if event.SubjectType != "incident" || event.SubjectID != incidentID.String() {
		t.Fatalf("unexpected incident audit subject: %+v", event)
	}
	if event.ClientID == nil || *event.ClientID != clientID.String() {
		t.Fatalf("incident audit client ID = %v, want %s", event.ClientID, clientID)
	}
	if event.AccessRule == nil || *event.AccessRule != "CLIENT.INCIDENT.VIEW" {
		t.Fatalf("incident audit access rule = %v", event.AccessRule)
	}
	if event.Details != nil {
		t.Fatalf("incident audit details = %#v, want no sensitive payload", event.Details)
	}
}
