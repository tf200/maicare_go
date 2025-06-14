package async

import (
	"context"
	"crypto/tls"
	"log"
	"maicare_go/util"
	"testing"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestEnqueueEmailDelivery(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}
	payload := EmailDeliveryPayload{
		To:           "farjiataha@gmail.com",
		UserEmail:    "farjiataha@gmail.com",
		UserPassword: "password",
	}
	ctx := context.Background()
	err = testasynqClient.EnqueueEmailDelivery(payload, ctx)
	t.Log(err)
	require.NoError(t, err)

	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:      config.RedisHost,
		Password:  config.RedisPassword,
		Username:  "",
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
	require.Equal(t, payload.To, gotPayload.To)

}
