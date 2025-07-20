package db

import (
	"context"
	"sync"
	"testing"

	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomSenders(t *testing.T) Sender {
	arg := CreateSenderParams{
		Types:        "main_provider",
		Name:         util.RandomString(5),
		Address:      util.StringPtr("test"),
		PostalCode:   util.StringPtr("test"),
		Place:        util.StringPtr("test"),
		Land:         util.StringPtr("test"),
		Kvknumber:    util.StringPtr("test"),
		Btwnumber:    util.StringPtr("test"),
		PhoneNumber:  util.StringPtr("test"),
		ClientNumber: util.StringPtr("test"),
		EmailAddress: util.StringPtr("test"),
		Contacts:     []byte(`[{"name": "Test Contact", "email": "test@example.com", "phone": "1234567890"}]`),
	}

	sender, err := testQueries.CreateSender(context.Background(), arg)
	require.NoError(t, err)
	templItems, err := testQueries.CreateSenderInvoiceTemplate(context.Background(), CreateSenderInvoiceTemplateParams{
		ID:              sender.ID,
		InvoiceTemplate: []int64{1, 2, 3},
	})

	require.NoError(t, err)
	require.NotEmpty(t, templItems)

	require.NotEmpty(t, sender)
	require.Equal(t, arg.Types, sender.Types)
	require.Equal(t, arg.Name, sender.Name)
	require.Equal(t, arg.Address, sender.Address)
	require.Equal(t, arg.PostalCode, sender.PostalCode)
	require.Equal(t, arg.Place, sender.Place)
	require.Equal(t, arg.Land, sender.Land)
	require.Equal(t, arg.Kvknumber, sender.Kvknumber)
	require.Equal(t, arg.Btwnumber, sender.Btwnumber)
	require.Equal(t, arg.PhoneNumber, sender.PhoneNumber)
	require.Equal(t, arg.ClientNumber, sender.ClientNumber)
	require.Equal(t, arg.EmailAddress, sender.EmailAddress)
	require.Equal(t, arg.Contacts, sender.Contacts)
	return sender
}

func TestCreateSender(t *testing.T) {
	createRandomSenders(t)
}

func TestListsenders(t *testing.T) {
	// Create test data
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomSenders(t)
		}()
	}
	wg.Wait()

	// Test multiple scenarios
	testCases := []struct {
		name  string
		arg   ListSendersParams
		check func(t *testing.T, senders []Sender)
	}{
		{
			name: "base case",
			arg: ListSendersParams{
				Limit:           5,
				Offset:          0,
				Search:          util.StringPtr(""),
				IncludeArchived: util.BoolPtr(true),
			},
			check: func(t *testing.T, senders []Sender) {
				require.Len(t, senders, 5)
			},
		},
		{
			name: "with offset",
			arg: ListSendersParams{
				Limit:           5,
				Offset:          5,
				Search:          util.StringPtr(""),
				IncludeArchived: util.BoolPtr(true),
			},
			check: func(t *testing.T, senders []Sender) {
				require.Len(t, senders, 5)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			senders, err := testQueries.ListSenders(context.Background(), tc.arg)
			require.NoError(t, err)
			tc.check(t, senders)
		})
	}
}

func TestCountSenders(t *testing.T) {
	// Get initial count before adding test data
	initialCount, err := testQueries.CountSenders(context.Background(), util.BoolPtr(true))
	require.NoError(t, err)

	// Create test data
	numSenders := 20
	var wg sync.WaitGroup
	for i := 0; i < numSenders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomSenders(t)
		}()
	}
	wg.Wait()

	testCases := []struct {
		name       string
		params     *bool
		checkCount func(t *testing.T, count int64)
	}{
		{
			name:   "Count all senders",
			params: util.BoolPtr(true),
			checkCount: func(t *testing.T, count int64) {
				require.Equal(t, initialCount+int64(numSenders), count,
					"should match initial count plus newly created employees")
			},
		},
		{
			name:   "Count non-archived senders",
			params: nil,
			checkCount: func(t *testing.T, count int64) {
				require.LessOrEqual(t, count, initialCount+int64(numSenders))
				require.GreaterOrEqual(t, count, initialCount)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := testQueries.CountSenders(context.Background(), tc.params)
			require.NoError(t, err)
			tc.checkCount(t, count)
		})
	}
}

func TestUpdateSender(t *testing.T) {
	sender := createRandomSenders(t)
	arg := UpdateSenderParams{
		ID:           sender.ID,
		Name:         util.StringPtr(util.RandomString(5)),
		Address:      util.StringPtr("test"),
		PostalCode:   util.StringPtr("test"),
		Place:        util.StringPtr("test"),
		Land:         util.StringPtr("test"),
		Kvknumber:    util.StringPtr("test"),
		Btwnumber:    util.StringPtr("test"),
		ClientNumber: util.StringPtr("test"),
		EmailAddress: util.StringPtr("test"),
		Contacts:     []byte(`[{"name": "Test Contact", "email": "example@gmail.com", "phone": "1234567890"}]`),
	}
	updatedSender, err := testQueries.UpdateSender(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSender)
	require.Equal(t, arg.ID, updatedSender.ID)
	require.NotNil(t, arg.Name)
	require.Equal(t, *arg.Name, updatedSender.Name)
	require.NotEqual(t, sender.Name, updatedSender.Name)
	require.NotEqual(t, sender.Contacts, updatedSender.Contacts)
}

func TestGetSenderById(t *testing.T) {
	sender := createRandomSenders(t)
	sender2, err := testQueries.GetSenderById(context.Background(), sender.ID)

	require.NoError(t, err)
	require.NotEmpty(t, sender2)
	require.Equal(t, sender.ID, sender2.ID)
	require.Equal(t, sender.Name, sender2.Name)
	require.Equal(t, sender.Address, sender2.Address)
}
