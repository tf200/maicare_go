package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestDischargeOverview(t *testing.T) {
	client := createRandomClientDetails(t)
	_ = createRandomContract(t, client.ID, client.SenderID)
	arg := CreateSchedueledClientStatusChangeParams{
		ClientID:      client.ID,
		NewStatus:     "Out Of Care",
		Reason:        util.StringPtr("Test Reason"),
		ScheduledDate: pgtype.Date{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	}

	_, err := testQueries.CreateSchedueledClientStatusChange(context.Background(), arg)
	require.NoError(t, err)

	overview, err := testQueries.DischargeOverview(context.Background(), DischargeOverviewParams{
		Limit:      5,
		Offset:     0,
		FilterType: "all",
	})
	require.NoError(t, err)
	require.NotEmpty(t, overview)

}

func TestTotalDischargeCount(t *testing.T) {
	client := createRandomClientDetails(t)
	_ = createRandomContract(t, client.ID, client.SenderID)
	arg := CreateSchedueledClientStatusChangeParams{
		ClientID:      client.ID,
		NewStatus:     "Out Of Care",
		Reason:        util.StringPtr("Test Reason"),
		ScheduledDate: pgtype.Date{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	}

	_, err := testQueries.CreateSchedueledClientStatusChange(context.Background(), arg)
	require.NoError(t, err)

	count, err := testQueries.TotalDischargeCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, count)
	t.Log(count)
}

func TestUrgentCasesCount(t *testing.T) {
	client1 := createRandomClientDetails(t)
	_ = createRandomContract(t, client1.ID, client1.SenderID)

	client2 := createRandomClientDetails(t)
	arg := CreateSchedueledClientStatusChangeParams{
		ClientID:      client2.ID,
		NewStatus:     "Out Of Care",
		Reason:        util.StringPtr("Test Reason"),
		ScheduledDate: pgtype.Date{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	}
	_, err := testQueries.CreateSchedueledClientStatusChange(context.Background(), arg)
	require.NoError(t, err)

	count, err := testQueries.UrgentCasesCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, count)
	t.Log(count)
}

func TestStatusChangeCount(t *testing.T) {
	client := createRandomClientDetails(t)
	arg := CreateSchedueledClientStatusChangeParams{
		ClientID:      client.ID,
		NewStatus:     "Out Of Care",
		Reason:        util.StringPtr("Test Reason"),
		ScheduledDate: pgtype.Date{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	}

	_, err := testQueries.CreateSchedueledClientStatusChange(context.Background(), arg)
	require.NoError(t, err)

	count, err := testQueries.StatusChangeCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, count)
}

func TestContractEndCount(t *testing.T) {
	client := createRandomClientDetails(t)
	_ = createRandomContract(t, client.ID, client.SenderID)

	count, err := testQueries.ContractEndCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, count)
}
