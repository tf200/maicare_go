package service

import (
	"context"
	"testing"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

func TestLogGroupEAuditExcludesClinicalPayload(t *testing.T) {
	recorder := &incidentAuditRecorder{}
	service := &ClientService{audit: recorder}
	clientID := uuid.New()
	evaluationID := uuid.New()

	service.logGroupEAudit(context.Background(), "read", "client_goal_evaluation", evaluationID, clientID, domain.PermClientEvaluationView.String(), 3)

	if len(recorder.events) != 1 {
		t.Fatalf("audit event count = %d, want 1", len(recorder.events))
	}
	event := recorder.events[0]
	if event.EventType != "record_access" || event.Action != "read" || event.Result != "success" {
		t.Fatalf("unexpected Group E audit classification: %+v", event)
	}
	if event.SubjectType != "client_goal_evaluation" || event.SubjectID != evaluationID.String() {
		t.Fatalf("unexpected Group E audit subject: %+v", event)
	}
	if event.ClientID == nil || *event.ClientID != clientID.String() {
		t.Fatalf("Group E audit client ID = %v, want %s", event.ClientID, clientID)
	}
	if event.AccessRule == nil || *event.AccessRule != domain.PermClientEvaluationView.String() {
		t.Fatalf("Group E audit access rule = %v", event.AccessRule)
	}
	if len(event.Details) != 1 || event.Details["count"] != 3 {
		t.Fatalf("Group E audit details = %#v, want count only", event.Details)
	}
}
