package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateAppointmentCard(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateAppointmentCardParams
		checks func(t *testing.T, card AppointmentCard, err error)
	}{
		{
			name: "successful appointment card creation with full fields",
			setup: func(ctx context.Context, qtx *Queries) CreateAppointmentCardParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateAppointmentCardParams{
					ClientID:               client.ID,
					GeneralInformation:     []string{"info1", "info2"},
					ImportantContacts:      []string{"contact1"},
					HouseholdInfo:          []string{"house1"},
					OrganizationAgreements: []string{"org1"},
					YouthOfficerAgreements: []string{"youth1"},
					TreatmentAgreements:    []string{"treat1"},
					SmokingRules:           []string{"smoke1"},
					Work:                   []string{"work1"},
					SchoolInternship:       []string{"school1"},
					Travel:                 []string{"travel1"},
					Leave:                  []string{"leave1"},
				}
			},
			checks: func(t *testing.T, card AppointmentCard, err error) {
				require.NoError(t, err, "CreateAppointmentCard should not return an error")
				require.NotZero(t, card.ID)
				require.Equal(t, 2, len(card.GeneralInformation))
			},
		},
		{
			name: "successful appointment card creation with empty arrays",
			setup: func(ctx context.Context, qtx *Queries) CreateAppointmentCardParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateAppointmentCardParams{
					ClientID:               client.ID,
					GeneralInformation:     []string{},
					ImportantContacts:      []string{},
					HouseholdInfo:          []string{},
					OrganizationAgreements: []string{},
					YouthOfficerAgreements: []string{},
					TreatmentAgreements:    []string{},
					SmokingRules:           []string{},
					Work:                   []string{},
					SchoolInternship:       []string{},
					Travel:                 []string{},
					Leave:                  []string{},
				}
			},
			checks: func(t *testing.T, card AppointmentCard, err error) {
				require.NoError(t, err, "CreateAppointmentCard should not return an error")
				require.NotZero(t, card.ID)
				require.Empty(t, card.GeneralInformation)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			card, err := qtx.CreateAppointmentCard(ctx, params)
			tt.checks(t, card, err)
		})
	}
}

func TestGetAppointmentCard(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, row GetAppointmentCardRow, err error)
	}{
		{
			name: "get existing appointment card",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomAppointmentCard(ctx, qtx, client.ID)
				return client.ID
			},
			checks: func(t *testing.T, row GetAppointmentCardRow, err error) {
				require.NoError(t, err, "GetAppointmentCard should not error")
				require.NotZero(t, row.ID)
			},
		},
		{
			name: "get appointment card for client without card",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				return client.ID
			},
			checks: func(t *testing.T, row GetAppointmentCardRow, err error) {
				require.Error(t, err, "GetAppointmentCard should error for client without card")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			clientID := tt.setup(ctx, qtx)
			row, err := qtx.GetAppointmentCard(ctx, clientID)
			tt.checks(t, row, err)
		})
	}
}

func TestUpdateAppointmentCard(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateAppointmentCardParams
		checks func(t *testing.T, card AppointmentCard, params UpdateAppointmentCardParams, err error)
	}{
		{
			name: "update existing appointment card",
			setup: func(ctx context.Context, qtx *Queries) UpdateAppointmentCardParams {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomAppointmentCard(ctx, qtx, client.ID)
				return UpdateAppointmentCardParams{
					ClientID:           client.ID,
					GeneralInformation: []string{"updated info"},
				}
			},
			checks: func(t *testing.T, card AppointmentCard, params UpdateAppointmentCardParams, err error) {
				require.NoError(t, err, "UpdateAppointmentCard should not error")
				require.Equal(t, params.GeneralInformation, card.GeneralInformation)
			},
		},
		{
			name: "update appointment card for client without card",
			setup: func(ctx context.Context, qtx *Queries) UpdateAppointmentCardParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateAppointmentCardParams{
					ClientID:           client.ID,
					GeneralInformation: []string{"info"},
				}
			},
			checks: func(t *testing.T, card AppointmentCard, params UpdateAppointmentCardParams, err error) {
				require.Error(t, err, "UpdateAppointmentCard should error for client without card")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			card, err := qtx.UpdateAppointmentCard(ctx, params)
			tt.checks(t, card, params, err)
		})
	}
}

func TestUpdateAppointmentCardUrl(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateAppointmentCardUrlParams
		checks func(t *testing.T, fileUrl *string, params UpdateAppointmentCardUrlParams, err error)
	}{
		{
			name: "update appointment card url",
			setup: func(ctx context.Context, qtx *Queries) UpdateAppointmentCardUrlParams {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomAppointmentCard(ctx, qtx, client.ID)
				url := util.RandomString(10)
				return UpdateAppointmentCardUrlParams{
					ClientID: client.ID,
					FileUrl:  &url,
				}
			},
			checks: func(t *testing.T, fileUrl *string, params UpdateAppointmentCardUrlParams, err error) {
				require.NoError(t, err, "UpdateAppointmentCardUrl should not error")
				require.Equal(t, params.FileUrl, fileUrl)
			},
		},
		{
			name: "update appointment card url for client without card",
			setup: func(ctx context.Context, qtx *Queries) UpdateAppointmentCardUrlParams {
				client := createRandomClientDetails(ctx, qtx)
				url := util.RandomString(10)
				return UpdateAppointmentCardUrlParams{
					ClientID: client.ID,
					FileUrl:  &url,
				}
			},
			checks: func(t *testing.T, fileUrl *string, params UpdateAppointmentCardUrlParams, err error) {
				require.Error(t, err, "UpdateAppointmentCardUrl should error for client without card")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			fileUrl, err := qtx.UpdateAppointmentCardUrl(ctx, params)
			tt.checks(t, fileUrl, params, err)
		})
	}
}

// Helpers
func createRandomAppointmentCard(ctx context.Context, qtx *Queries, clientID uuid.UUID) AppointmentCard {

	params := CreateAppointmentCardParams{
		ClientID:               clientID,
		GeneralInformation:     []string{util.RandomString(5)},
		ImportantContacts:      []string{util.RandomString(5)},
		HouseholdInfo:          []string{util.RandomString(5)},
		OrganizationAgreements: []string{util.RandomString(5)},
		YouthOfficerAgreements: []string{util.RandomString(5)},
		TreatmentAgreements:    []string{util.RandomString(5)},
		SmokingRules:           []string{util.RandomString(5)},
		Work:                   []string{util.RandomString(5)},
		SchoolInternship:       []string{util.RandomString(5)},
		Travel:                 []string{util.RandomString(5)},
		Leave:                  []string{util.RandomString(5)},
	}
	card, err := qtx.CreateAppointmentCard(ctx, params)
	if err != nil {
		panic(err)
	}
	return card
}
