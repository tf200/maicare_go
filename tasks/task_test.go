package tasks

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"maicare_go/util"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestEnqueueEmailDelivery(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}
	payload := EmailDeliveryPayload{
		UserID:     1,
		TemplateID: "template_id",
	}
	ctx := context.Background()
	err = testasynqClient.EnqueueEmailDelivery(payload, ctx)
	t.Log(err)
	require.NoError(t, err)

	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:      config.RedisHost,
		Password:  config.RedisPassword,
		Username:  config.RedisUser,
		TLSConfig: &tls.Config{},
	})

	tasks, err := inspector.ListPendingTasks("default")
	require.NoError(t, err)
	require.NotEmpty(t, tasks)

	task := tasks[0]
	require.Equal(t, TypeEmailDelivery, task.Type)

	// Decode and verify payload
	var gotPayload EmailDeliveryPayload
	err = json.Unmarshal(task.Payload, &gotPayload)
	require.NoError(t, err)
	require.Equal(t, payload.UserID, gotPayload.UserID)

}
