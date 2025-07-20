package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchInvoiceTemplateItems(t *testing.T) {
	ctx := context.Background()
	store := NewStore(testDB)
	client := createRandomClientDetails(t)
	sender := createRandomSenders(t)

	contract := createRandomContract(t, client.ID, &sender.ID)

	data := FetchQueryData{
		ClientID:   client.ID,
		ContractID: contract.ID,
		SenderID:   sender.ID,
	}

	extraContent, err := store.FetchInvoiceTemplateItems(ctx, data)
	t.Log(err)
	require.NoError(t, err)
	require.NotEmpty(t, extraContent)
}
