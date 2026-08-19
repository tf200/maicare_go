package worker

import (
	"testing"

	"maicare_go/internal/ctxkeys"
	pkgasynq "maicare_go/pkg/asynq"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

func TestSchedulerTaskIncludesSystemActor(t *testing.T) {
	actor := ctxkeys.ActorIdentity{UserID: uuid.New(), EmployeeID: uuid.New()}
	scheduler := &Scheduler{actor: actor}

	task, err := scheduler.task(TypeContractReminder)
	if err != nil {
		t.Fatalf("task() error = %v", err)
	}
	var payload pkgasynq.ScheduledWorkerPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload.Actor.UserID != actor.UserID || payload.Actor.EmployeeID != actor.EmployeeID {
		t.Fatalf("actor = %#v, want %#v", payload.Actor, actor)
	}
}

func TestSchedulerTaskRejectsMissingSystemActor(t *testing.T) {
	scheduler := &Scheduler{}
	if _, err := scheduler.task(TypeContractReminder); err == nil {
		t.Fatal("expected missing system actor error")
	}
}
