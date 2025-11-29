package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateSender(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateSenderParams
		checks func(t *testing.T, sender Sender, params CreateSenderParams)
	}{
		{
			name: "Create Sender Successfully with minial fields",
			setup: func(ctx context.Context, qtx *Queries) CreateSenderParams {
				return CreateSenderParams{
					Types: SenderTypesEnumParticularParty,
					Name:  util.RandomString(5),
				}
			},
			checks: func(t *testing.T, sender Sender, params CreateSenderParams) {
				require.Equal(t, params.Types, sender.Types)
				require.Equal(t, params.Name, sender.Name)
				require.Nil(t, sender.Address)
				require.Nil(t, sender.PostalCode)
				require.Nil(t, sender.Place)
				require.Nil(t, sender.Land)
				require.Nil(t, sender.Kvknumber)
				require.Nil(t, sender.Btwnumber)
				require.Nil(t, sender.PhoneNumber)
				require.Nil(t, sender.ClientNumber)
				require.Nil(t, sender.EmailAddress)
				require.Empty(t, sender.Contacts)
			},
		},
		{
			name: "Create Sender Successfully with all fields",
			setup: func(ctx context.Context, qtx *Queries) CreateSenderParams {
				return CreateSenderParams{
					Types:        SenderTypesEnumMainProvider,
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
			},
			checks: func(t *testing.T, sender Sender, params CreateSenderParams) {
				require.Equal(t, params.Types, sender.Types)
				require.Equal(t, params.Name, sender.Name)
				require.Equal(t, params.Address, sender.Address)
				require.Equal(t, params.PostalCode, sender.PostalCode)
				require.Equal(t, params.Place, sender.Place)
				require.Equal(t, params.Land, sender.Land)
				require.Equal(t, params.Kvknumber, sender.Kvknumber)
				require.Equal(t, params.Btwnumber, sender.Btwnumber)
				require.Equal(t, params.PhoneNumber, sender.PhoneNumber)
				require.Equal(t, params.ClientNumber, sender.ClientNumber)
				require.Equal(t, params.EmailAddress, sender.EmailAddress)
				require.Equal(t, params.Contacts, sender.Contacts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			sender, err := qtx.CreateSender(ctx, params)
			require.NoError(t, err)
			require.NotEmpty(t, sender)

			tt.checks(t, sender, params)
		})
	}
}
func TestGetSenderById(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, sender Sender, err error)
	}{
		{
			name: "Get existing sender by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				sender := createRandomSenders(ctx, qtx)
				return sender.ID
			},
			checks: func(t *testing.T, sender Sender, err error) {
				require.NoError(t, err, "GetSenderById() should not error")
				require.NotEmpty(t, sender.ID)
				require.Equal(t, SenderTypesEnumHealthcareInstitution, sender.Types)
				require.NotEmpty(t, sender.Name)
			},
		},
		{
			name: "Get non-existent sender by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New() // Assuming this ID doesn't exist
			},
			checks: func(t *testing.T, sender Sender, err error) {
				require.Error(t, err, "GetSenderById() should error for non-existent ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			sender, err := qtx.GetSenderById(ctx, id)
			tt.checks(t, sender, err)
		})
	}
}

func TestUpdateSender(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateSenderParams
		checks func(t *testing.T, sender Sender, err error)
	}{
		{
			name: "Update sender successfully",
			setup: func(ctx context.Context, qtx *Queries) UpdateSenderParams {
				sender := createRandomSenders(ctx, qtx)
				newName := util.RandomString(5)
				return UpdateSenderParams{
					ID:   sender.ID,
					Name: &newName,
				}
			},
			checks: func(t *testing.T, sender Sender, err error) {
				require.NoError(t, err, "UpdateSender() should not error")
				require.NotEmpty(t, sender.Name)
			},
		},
		{
			name: "Update non-existent sender",
			setup: func(ctx context.Context, qtx *Queries) UpdateSenderParams {
				newName := util.RandomString(5)
				return UpdateSenderParams{
					ID:   uuid.New(),
					Name: &newName,
				}
			},
			checks: func(t *testing.T, sender Sender, err error) {
				require.NoError(t, err, "UpdateSender() should not error even for non-existent sender")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			sender, err := qtx.UpdateSender(ctx, params)
			tt.checks(t, sender, err)
		})
	}
}

func TestDeleteSender(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "Delete existing sender",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				sender := createRandomSenders(ctx, qtx)
				return sender.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteSender() should not error")
			},
		},
		{
			name: "Delete non-existent sender",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteSender() should not error for non-existent sender")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			err = qtx.DeleteSender(ctx, id)
			tt.checks(t, err)
		})
	}
}

func TestListSenders(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListSendersParams
		checks func(t *testing.T, senders []Sender, err error)
	}{
		{
			name: "List senders without filters",
			setup: func(ctx context.Context, qtx *Queries) ListSendersParams {
				createRandomSenders(ctx, qtx)
				return ListSendersParams{
					Limit:  10,
					Offset: 0,
				}
			},
			checks: func(t *testing.T, senders []Sender, err error) {
				require.NoError(t, err, "ListSenders() should not error")
				require.GreaterOrEqual(t, len(senders), 1)
			},
		},
		{
			name: "List senders with search",
			setup: func(ctx context.Context, qtx *Queries) ListSendersParams {
				sender := createRandomSenders(ctx, qtx)
				search := sender.Name[:3] // Partial name
				return ListSendersParams{
					Limit:  10,
					Offset: 0,
					Search: &search,
				}
			},
			checks: func(t *testing.T, senders []Sender, err error) {
				require.NoError(t, err, "ListSenders() should not error")
				require.GreaterOrEqual(t, len(senders), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			senders, err := qtx.ListSenders(ctx, params)
			tt.checks(t, senders, err)
		})
	}
}

func TestCountSenders(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) *bool
		checks func(t *testing.T, count int64, err error)
	}{
		{
			name: "Count active senders",
			setup: func(ctx context.Context, qtx *Queries) *bool {
				createRandomSenders(ctx, qtx)
				includeArchived := false
				return &includeArchived
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountSenders() should not error")
				require.GreaterOrEqual(t, count, 1)
			},
		},
		{
			name: "Count all senders including archived",
			setup: func(ctx context.Context, qtx *Queries) *bool {
				createRandomSenders(ctx, qtx)
				includeArchived := true
				return &includeArchived
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountSenders() should not error")
				require.GreaterOrEqual(t, count, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			includeArchived := tt.setup(ctx, qtx)
			count, err := qtx.CountSenders(ctx, includeArchived)
			tt.checks(t, count, err)
		})
	}
}

func TestCreateSenderInvoiceTemplate(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateSenderInvoiceTemplateParams
		checks func(t *testing.T, template []uuid.UUID, err error)
	}{
		{
			name: "Create invoice template successfully",
			setup: func(ctx context.Context, qtx *Queries) CreateSenderInvoiceTemplateParams {
				sender := createRandomSenders(ctx, qtx)
				// list templateItems uuids from database
				templates, err := qtx.GetAllTemplateItems(ctx)
				require.NoError(t, err, "GetAllTemplateItems() should not error")
				return CreateSenderInvoiceTemplateParams{
					ID:              sender.ID,
					InvoiceTemplate: []uuid.UUID{templates[0].ID, templates[1].ID, templates[2].ID},
				}
			},
			checks: func(t *testing.T, template []uuid.UUID, err error) {
				require.NoError(t, err, "CreateSenderInvoiceTemplate() should not error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			template, err := qtx.CreateSenderInvoiceTemplate(ctx, params)
			tt.checks(t, template, err)
		})
	}
}

func TestGetSenderInvoiceTemplate(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, template []uuid.UUID, err error)
	}{
		{
			name: "Get invoice template for existing sender",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				sender := createRandomSenders(ctx, qtx)
				return sender.ID
			},
			checks: func(t *testing.T, template []uuid.UUID, err error) {
				require.NoError(t, err, "GetSenderInvoiceTemplate() should not error")
			},
		},
		{
			name: "Get invoice template for non-existent sender",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, template []uuid.UUID, err error) {
				require.Error(t, err, "GetSenderInvoiceTemplate() should error for non-existent sender")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := testDB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)
			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			template, err := qtx.GetSenderInvoiceTemplate(ctx, id)
			tt.checks(t, template, err)
		})
	}
}
func createRandomSenders(ctx context.Context, qtx *Queries) Sender {
	arg := CreateSenderParams{
		Types:        SenderTypesEnumHealthcareInstitution,
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

	sender, err := qtx.CreateSender(context.Background(), arg)
	if err != nil {
		panic(err)
	}
	templates, err := qtx.GetAllTemplateItems(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = qtx.CreateSenderInvoiceTemplate(context.Background(), CreateSenderInvoiceTemplateParams{
		ID:              sender.ID,
		InvoiceTemplate: []uuid.UUID{templates[0].ID, templates[1].ID, templates[2].ID},
	})
	if err != nil {
		panic(err)
	}

	return sender
}
