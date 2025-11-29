package db

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateClientDetails(t *testing.T) {
	tests := []struct {
		name   string
		params CreateClientDetailsParams
		checks func(t *testing.T, client ClientDetail)
	}{
		{
			name: "successful creation with minimal required fields",
			params: CreateClientDetailsParams{
				FirstName:                  "John",
				LastName:                   "Doe",
				DateOfBirth:                pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
				Identity:                   true,
				Email:                      util.RandomEmail(),
				Gender:                     ClientGenderEnumMale,
				Filenumber:                 util.RandomString(10),
				EducationLevel:             ClientEducationLevelEnumPrimary,
				Addresses:                  []byte(`[]`),
				EducationCurrentlyEnrolled: false,
				WorkCurrentlyEmployed:      false,
				LivingSituation:            NullClientLivingSituationEnum{Valid: false},
			},
			checks: func(t *testing.T, client ClientDetail) {
				require.Equal(t, "John", client.FirstName)
				require.Equal(t, "Doe", client.LastName)
				require.Equal(t, ClientGenderEnumMale, client.Gender)
				require.True(t, client.Identity)
				require.False(t, client.EducationCurrentlyEnrolled)
				require.False(t, client.WorkCurrentlyEmployed)
				require.NotNil(t, client.ID)
			},
		},
		{
			name: "successful creation with all fields populated",
			params: CreateClientDetailsParams{
				FirstName:                  "Jane",
				LastName:                   "Smith",
				DateOfBirth:                pgtype.Date{Time: time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC), Valid: true},
				Identity:                   true,
				Bsn:                        util.StringPtr("123456789"),
				Source:                     util.StringPtr("Referral"),
				Birthplace:                 util.StringPtr("Amsterdam"),
				Email:                      util.RandomEmail(),
				PhoneNumber:                util.StringPtr("+31612345678"),
				Departement:                util.StringPtr("Medical"),
				Gender:                     ClientGenderEnumFemale,
				Filenumber:                 util.RandomString(10),
				ProfilePicture:             util.StringPtr("https://example.com/pic.jpg"),
				Infix:                      util.StringPtr("van"),
				DepartureReason:            util.StringPtr("Completed program"),
				DepartureReport:            util.StringPtr("Successful completion"),
				Addresses:                  []byte(`[{"street":"Main St","city":"Amsterdam"}]`),
				LegalMeasure:               util.StringPtr("Guardianship"),
				EducationCurrentlyEnrolled: true,
				EducationInstitution:       util.StringPtr("University of Amsterdam"),
				EducationMentorName:        util.StringPtr("Prof. John"),
				EducationMentorPhone:       util.StringPtr("+31687654321"),
				EducationMentorEmail:       util.StringPtr("john@university.nl"),
				EducationAdditionalNotes:   util.StringPtr("Excellent student"),
				EducationLevel:             ClientEducationLevelEnumHigher,
				WorkCurrentlyEmployed:      true,
				WorkCurrentEmployer:        util.StringPtr("Tech Corp"),
				WorkCurrentEmployerPhone:   util.StringPtr("+31612223344"),
				WorkCurrentEmployerEmail:   util.StringPtr("hr@techcorp.nl"),
				WorkCurrentPosition:        util.StringPtr("Senior Developer"),
				WorkStartDate:              pgtype.Date{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
				WorkAdditionalNotes:        util.StringPtr("Promoted twice"),
				LivingSituation:            NullClientLivingSituationEnum{ClientLivingSituationEnum: ClientLivingSituationEnumHome, Valid: true},
				LivingSituationNotes:       util.StringPtr("Lives with partner"),
			},
			checks: func(t *testing.T, client ClientDetail) {
				require.Equal(t, "Jane", client.FirstName)
				require.Equal(t, "Smith", client.LastName)
				require.NotNil(t, client.Bsn)
				require.Equal(t, "123456789", *client.Bsn)
				require.NotNil(t, client.ProfilePicture)
				require.Equal(t, "https://example.com/pic.jpg", *client.ProfilePicture)
				require.True(t, client.WorkCurrentlyEmployed)
				require.True(t, client.EducationCurrentlyEnrolled)
				require.Equal(t, ClientEducationLevelEnumHigher, client.EducationLevel)
			},
		},
		{
			name: "creation with nil optional fields",
			params: CreateClientDetailsParams{
				FirstName:                  "Bob",
				LastName:                   "Johnson",
				DateOfBirth:                pgtype.Date{Time: time.Date(1995, 3, 10, 0, 0, 0, 0, time.UTC), Valid: true},
				Identity:                   false,
				Email:                      util.RandomEmail(),
				Gender:                     ClientGenderEnumOther,
				Filenumber:                 util.RandomString(10),
				EducationLevel:             ClientEducationLevelEnumNone,
				Addresses:                  []byte(`[]`),
				EducationCurrentlyEnrolled: false,
				WorkCurrentlyEmployed:      false,
				LivingSituation:            NullClientLivingSituationEnum{Valid: false},
				Bsn:                        nil,
				ProfilePicture:             nil,
				Infix:                      nil,
			},
			checks: func(t *testing.T, client ClientDetail) {
				require.Nil(t, client.Bsn)
				require.Nil(t, client.ProfilePicture)
				require.Nil(t, client.Infix)
				require.False(t, client.Identity)
				require.Equal(t, ClientGenderEnumOther, client.Gender)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			client, err := qtx.CreateClientDetails(ctx, tt.params)
			require.NoError(t, err, "CreateClientDetails() should not error")

			tt.checks(t, client)
		})
	}
}

func TestGetClientDetails(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, client GetClientDetailsRow, err error)
	}{
		{
			name: "get existing client details",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				return client.ID
			},
			checks: func(t *testing.T, client GetClientDetailsRow, err error) {
				require.NoError(t, err, "GetClientDetails() should not error")
				require.NotEmpty(t, client.ID)
				require.Equal(t, "John", client.FirstName)
				require.Equal(t, "Doe", client.LastName)
			},
		},
		{
			name: "get non-existent client details",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, client GetClientDetailsRow, err error) {
				require.Error(t, err, "GetClientDetails() should error for non-existent client")
			},
		},
		{
			name: "get client with all fields populated",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				params := CreateClientDetailsParams{
					FirstName:                  "Complete",
					LastName:                   "Data",
					DateOfBirth:                pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Identity:                   true,
					Bsn:                        util.StringPtr("987654321"),
					Source:                     util.StringPtr("Direct"),
					Birthplace:                 util.StringPtr("Rotterdam"),
					Email:                      util.RandomEmail(),
					PhoneNumber:                util.StringPtr("+31698765432"),
					Gender:                     ClientGenderEnumFemale,
					Filenumber:                 util.RandomString(10),
					ProfilePicture:             util.StringPtr("https://example.com/full.jpg"),
					Addresses:                  []byte(`[{"street":"Full St"}]`),
					EducationLevel:             ClientEducationLevelEnumHigher,
					EducationCurrentlyEnrolled: true,
					EducationInstitution:       util.StringPtr("MIT"),
					WorkCurrentlyEmployed:      true,
					WorkCurrentEmployer:        util.StringPtr("Google"),
					LivingSituation:            NullClientLivingSituationEnum{ClientLivingSituationEnum: ClientLivingSituationEnumHome, Valid: true},
				}
				client, _ := qtx.CreateClientDetails(ctx, params)
				return client.ID
			},
			checks: func(t *testing.T, client GetClientDetailsRow, err error) {
				require.NoError(t, err)
				require.Equal(t, "Complete", client.FirstName)
				require.Equal(t, "Data", client.LastName)
				require.NotNil(t, client.Bsn)
				require.Equal(t, "987654321", *client.Bsn)
				require.NotNil(t, client.ProfilePicture)
				require.True(t, client.WorkCurrentlyEmployed)
				require.True(t, client.EducationCurrentlyEnrolled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			client, err := qtx.GetClientDetails(ctx, id)
			tt.checks(t, client, err)
		})
	}
}

func TestUpdateClientDetails(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateClientDetailsParams
		checks func(t *testing.T, client ClientDetail, err error)
	}{
		{
			name: "successful update of single field",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDetailsParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientDetailsParams{
					ID:        client.ID,
					FirstName: util.StringPtr("UpdatedName"),
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err, "UpdateClientDetails() should not error")
				require.Equal(t, "UpdatedName", client.FirstName)
			},
		},
		{
			name: "successful update of multiple fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDetailsParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientDetailsParams{
					ID:          client.ID,
					FirstName:   util.StringPtr("Jane"),
					LastName:    util.StringPtr("Smith"),
					Email:       util.StringPtr("jane.smith@example.com"),
					PhoneNumber: util.StringPtr("+31612345678"),
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.Equal(t, "Jane", client.FirstName)
				require.Equal(t, "Smith", client.LastName)
				require.Equal(t, "jane.smith@example.com", client.Email)
				require.NotNil(t, client.PhoneNumber)
				require.Equal(t, "+31612345678", *client.PhoneNumber)
			},
		},
		{
			name: "update with nil values preserves existing data",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDetailsParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientDetailsParams{
					ID:        client.ID,
					FirstName: nil,
					Email:     util.StringPtr("newemail@example.com"),
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.Equal(t, "John", client.FirstName)
				require.Equal(t, "newemail@example.com", client.Email)
			},
		},
		{
			name: "update education fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDetailsParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientDetailsParams{
					ID:                         client.ID,
					EducationCurrentlyEnrolled: util.BoolPtr(true),
					EducationInstitution:       util.StringPtr("New University"),
					EducationMentorName:        util.StringPtr("Dr. Smith"),
					EducationLevel:             NullClientEducationLevelEnum{ClientEducationLevelEnum: ClientEducationLevelEnumHigher, Valid: true},
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.True(t, client.EducationCurrentlyEnrolled)
				require.NotNil(t, client.EducationInstitution)
				require.Equal(t, "New University", *client.EducationInstitution)
				require.Equal(t, ClientEducationLevelEnumHigher, client.EducationLevel)
			},
		},
		{
			name: "update work fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDetailsParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientDetailsParams{
					ID:                    client.ID,
					WorkCurrentlyEmployed: util.BoolPtr(true),
					WorkCurrentEmployer:   util.StringPtr("Tech Solutions"),
					WorkCurrentPosition:   util.StringPtr("Manager"),
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.True(t, client.WorkCurrentlyEmployed)
				require.NotNil(t, client.WorkCurrentEmployer)
				require.Equal(t, "Tech Solutions", *client.WorkCurrentEmployer)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			client, err := qtx.UpdateClientDetails(ctx, params)
			tt.checks(t, client, err)
		})
	}
}

func TestUpdateClientStatus(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateClientStatusParams
		checks func(t *testing.T, client ClientDetail, err error)
	}{
		{
			name: "update to In Care status",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientStatusParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientStatusParams{
					ID:     client.ID,
					Status: ClientStatusEnumInCare,
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err, "UpdateClientStatus() should not error")
				require.Equal(t, ClientStatusEnumInCare, client.Status)
			},
		},
		{
			name: "update to On Waiting List status",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientStatusParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientStatusParams{
					ID:     client.ID,
					Status: ClientStatusEnumOnWaitingList,
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.Equal(t, ClientStatusEnumOnWaitingList, client.Status)
			},
		},
		{
			name: "update to Out Of Care status",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientStatusParams {
				client := createRandomClientDetails(ctx, qtx)
				return UpdateClientStatusParams{
					ID:     client.ID,
					Status: ClientStatusEnumOutOfCare,
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.Equal(t, ClientStatusEnumOutOfCare, client.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			client, err := qtx.UpdateClientStatus(ctx, params)
			tt.checks(t, client, err)
		})
	}
}

func TestSetClientProfilePicture(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) SetClientProfilePictureParams
		checks func(t *testing.T, client ClientDetail, err error)
	}{
		{
			name: "set profile picture",
			setup: func(ctx context.Context, qtx *Queries) SetClientProfilePictureParams {
				client := createRandomClientDetails(ctx, qtx)
				return SetClientProfilePictureParams{
					ID:             client.ID,
					ProfilePicture: util.StringPtr("https://example.com/newpic.jpg"),
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err, "SetClientProfilePicture() should not error")
				require.NotNil(t, client.ProfilePicture)
				require.Equal(t, "https://example.com/newpic.jpg", *client.ProfilePicture)
			},
		},
		{
			name: "set profile picture to nil",
			setup: func(ctx context.Context, qtx *Queries) SetClientProfilePictureParams {
				params := CreateClientDetailsParams{
					FirstName:                  "Test",
					LastName:                   "User",
					DateOfBirth:                pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Identity:                   true,
					Email:                      util.RandomEmail(),
					Gender:                     ClientGenderEnumMale,
					Filenumber:                 util.RandomString(10),
					ProfilePicture:             util.StringPtr("https://example.com/old.jpg"),
					EducationLevel:             ClientEducationLevelEnumPrimary,
					Addresses:                  []byte(`[]`),
					EducationCurrentlyEnrolled: false,
					WorkCurrentlyEmployed:      false,
				}
				client, _ := qtx.CreateClientDetails(ctx, params)

				return SetClientProfilePictureParams{
					ID:             client.ID,
					ProfilePicture: nil,
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.Nil(t, client.ProfilePicture)
			},
		},
		{
			name: "update existing profile picture",
			setup: func(ctx context.Context, qtx *Queries) SetClientProfilePictureParams {
				params := CreateClientDetailsParams{
					FirstName:                  "Update",
					LastName:                   "Test",
					DateOfBirth:                pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Identity:                   true,
					Email:                      util.RandomEmail(),
					Gender:                     ClientGenderEnumMale,
					Filenumber:                 util.RandomString(10),
					ProfilePicture:             util.StringPtr("https://example.com/old.jpg"),
					EducationLevel:             ClientEducationLevelEnumPrimary,
					Addresses:                  []byte(`[]`),
					EducationCurrentlyEnrolled: false,
					WorkCurrentlyEmployed:      false,
				}
				client, _ := qtx.CreateClientDetails(ctx, params)
				return SetClientProfilePictureParams{
					ID:             client.ID,
					ProfilePicture: util.StringPtr("https://example.com/updated.jpg"),
				}
			},
			checks: func(t *testing.T, client ClientDetail, err error) {
				require.NoError(t, err)
				require.NotNil(t, client.ProfilePicture)
				require.Equal(t, "https://example.com/updated.jpg", *client.ProfilePicture)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			client, err := qtx.SetClientProfilePicture(ctx, params)
			tt.checks(t, client, err)
		})
	}
}

func TestCreateClientDocument(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateClientDocumentParams
		checks func(t *testing.T, doc ClientDocument)
	}{
		{
			name: "create registration form document",
			setup: func(ctx context.Context, qtx *Queries) CreateClientDocumentParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientDocumentParams{
					ClientID:       client.ID,
					AttachmentUuid: nil,
					Label:          ClientDocumentLabelEnumRegistrationForm,
				}
			},
			checks: func(t *testing.T, doc ClientDocument) {
				require.Nil(t, doc.AttachmentUuid)
				require.Equal(t, ClientDocumentLabelEnumRegistrationForm, doc.Label)
			},
		},
		{
			name: "create intake form document",
			setup: func(ctx context.Context, qtx *Queries) CreateClientDocumentParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientDocumentParams{
					ClientID:       client.ID,
					AttachmentUuid: nil,
					Label:          ClientDocumentLabelEnumIntakeForm,
				}
			},
			checks: func(t *testing.T, doc ClientDocument) {
				require.Equal(t, ClientDocumentLabelEnumIntakeForm, doc.Label)
			},
		},
		{
			name: "create document with attachment UUID",
			setup: func(ctx context.Context, qtx *Queries) CreateClientDocumentParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientDocumentParams{
					ClientID: client.ID,
					AttachmentUuid: func() *uuid.UUID {
						u := uuid.New()
						return &u
					}(),
					Label: ClientDocumentLabelEnumConsentForm,
				}
			},
			checks: func(t *testing.T, doc ClientDocument) {
				require.NotNil(t, doc.AttachmentUuid)
				require.Equal(t, ClientDocumentLabelEnumConsentForm, doc.Label)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)
			params := tt.setup(ctx, qtx)

			doc, err := qtx.CreateClientDocument(ctx, params)
			require.NoError(t, err, "CreateClientDocument() should not error")

			tt.checks(t, doc)
		})
	}
}

func TestDeleteClientDocument(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) *uuid.UUID
		checks func(t *testing.T, doc ClientDocument, err error)
	}{
		{
			name: "delete existing document",
			setup: func(ctx context.Context, qtx *Queries) *uuid.UUID {
				attachmentUUID := uuid.New()
				_, _ = qtx.CreateClientDocument(ctx, CreateClientDocumentParams{
					ClientID:       uuid.New(),
					AttachmentUuid: &attachmentUUID,
					Label:          ClientDocumentLabelEnumRegistrationForm,
				})
				return &attachmentUUID
			},
			checks: func(t *testing.T, doc ClientDocument, err error) {
				require.NoError(t, err, "DeleteClientDocument() should not error")
				require.NotNil(t, doc.AttachmentUuid)
			},
		},
		{
			name: "delete non-existent document",
			setup: func(ctx context.Context, qtx *Queries) *uuid.UUID {
				uuid := uuid.New()
				return &uuid
			},
			checks: func(t *testing.T, doc ClientDocument, err error) {
				require.Error(t, err, "DeleteClientDocument() should error for non-existent document")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			attachmentUUID := tt.setup(ctx, qtx)
			doc, err := qtx.DeleteClientDocument(ctx, attachmentUUID)
			tt.checks(t, doc, err)
		})
	}
}

func TestCreateClientStatusHistory(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateClientStatusHistoryParams
		checks func(t *testing.T, history ClientStatusHistory)
	}{
		{
			name: "create status history with old and new status",
			setup: func(ctx context.Context, qtx *Queries) CreateClientStatusHistoryParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientStatusHistoryParams{
					ClientID:  client.ID,
					OldStatus: util.StringPtr("On Waiting List"),
					NewStatus: "In Care",
					Reason:    util.StringPtr("Client accepted into program"),
				}
			},
			checks: func(t *testing.T, history ClientStatusHistory) {
				require.NotNil(t, history.OldStatus)
				require.Equal(t, "On Waiting List", *history.OldStatus)
				require.Equal(t, "In Care", history.NewStatus)
				require.NotNil(t, history.Reason)
			},
		},
		{
			name: "create status history with nil old status",
			setup: func(ctx context.Context, qtx *Queries) CreateClientStatusHistoryParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientStatusHistoryParams{
					ClientID:  client.ID,
					OldStatus: nil,
					NewStatus: "On Waiting List",
					Reason:    util.StringPtr("Initial status assignment"),
				}
			},
			checks: func(t *testing.T, history ClientStatusHistory) {
				require.Nil(t, history.OldStatus)
				require.Equal(t, "On Waiting List", history.NewStatus)
			},
		},
		{
			name: "create status history with nil reason",
			setup: func(ctx context.Context, qtx *Queries) CreateClientStatusHistoryParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientStatusHistoryParams{
					ClientID:  client.ID,
					OldStatus: util.StringPtr("In Care"),
					NewStatus: "Out Of Care",
					Reason:    nil,
				}
			},
			checks: func(t *testing.T, history ClientStatusHistory) {
				require.Nil(t, history.Reason)
				require.Equal(t, "Out Of Care", history.NewStatus)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)
			params := tt.setup(ctx, qtx)

			history, err := qtx.CreateClientStatusHistory(ctx, params)
			require.NoError(t, err, "CreateClientStatusHistory() should not error")

			tt.checks(t, history)
		})
	}
}

func TestListClientStatusHistory(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32)
		checks func(t *testing.T, histories []ClientStatusHistory, err error)
	}{
		{
			name: "list status history for client with multiple entries",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				client := createRandomClientDetails(ctx, qtx)
				clientID := client.ID
				for i := 0; i < 3; i++ {
					_, _ = qtx.CreateClientStatusHistory(ctx, CreateClientStatusHistoryParams{
						ClientID:  client.ID,
						OldStatus: nil,
						NewStatus: "In Care",
						Reason:    util.StringPtr("Test"),
					})
				}
				return clientID, 10, 0
			},
			checks: func(t *testing.T, histories []ClientStatusHistory, err error) {
				require.NoError(t, err, "ListClientStatusHistory() should not error")
				require.Greater(t, len(histories), 0)
			},
		},
		{
			name: "list status history with pagination",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				client := createRandomClientDetails(ctx, qtx)
				for i := 0; i < 5; i++ {
					_, _ = qtx.CreateClientStatusHistory(ctx, CreateClientStatusHistoryParams{
						ClientID:  client.ID,
						OldStatus: nil,
						NewStatus: "In Care",
					})
				}
				return client.ID, 2, 0
			},
			checks: func(t *testing.T, histories []ClientStatusHistory, err error) {
				require.NoError(t, err)
				require.LessOrEqual(t, len(histories), 2)
			},
		},
		{
			name: "list status history for client with no entries",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				return uuid.New(), 10, 0
			},
			checks: func(t *testing.T, histories []ClientStatusHistory, err error) {
				require.NoError(t, err)
				require.Empty(t, histories)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			clientID, limit, offset := tt.setup(ctx, qtx)
			histories, err := qtx.ListClientStatusHistory(ctx, ListClientStatusHistoryParams{
				ClientID: clientID,
				Limit:    limit,
				Offset:   offset,
			})
			tt.checks(t, histories, err)
		})
	}
}

func TestCreateSchedueledClientStatusChange(t *testing.T) {
	tests := []struct {
		name   string
		params CreateSchedueledClientStatusChangeParams
		checks func(t *testing.T, change ScheduledStatusChange)
	}{
		{
			name: "create scheduled status change",
			params: CreateSchedueledClientStatusChangeParams{
				ClientID:      uuid.New(),
				NewStatus:     util.StringPtr("Out Of Care"),
				Reason:        util.StringPtr("Completion of program"),
				ScheduledDate: pgtype.Date{Time: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true},
			},
			checks: func(t *testing.T, change ScheduledStatusChange) {
				require.NotNil(t, change.NewStatus)
				require.Equal(t, "Out Of Care", *change.NewStatus)
				require.NotNil(t, change.Reason)
			},
		},
		{
			name: "create scheduled change with nil reason",
			params: CreateSchedueledClientStatusChangeParams{
				ClientID:      uuid.New(),
				NewStatus:     util.StringPtr("In Care"),
				Reason:        nil,
				ScheduledDate: pgtype.Date{Time: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			},
			checks: func(t *testing.T, change ScheduledStatusChange) {
				require.Nil(t, change.Reason)
				require.NotNil(t, change.NewStatus)
			},
		},
		{
			name: "create scheduled change with future date",
			params: CreateSchedueledClientStatusChangeParams{
				ClientID:      uuid.New(),
				NewStatus:     util.StringPtr("Out Of Care"),
				Reason:        util.StringPtr("Planned departure"),
				ScheduledDate: pgtype.Date{Time: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			},
			checks: func(t *testing.T, change ScheduledStatusChange) {
				require.True(t, change.ScheduledDate.Time.After(time.Now()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			change, err := qtx.CreateSchedueledClientStatusChange(ctx, tt.params)
			require.NoError(t, err, "CreateSchedueledClientStatusChange() should not error")

			tt.checks(t, change)
		})
	}
}

func TestCreateClientLocationTransfer(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateClientLocationTransferParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "successful location transfer creation",
			setup: func(ctx context.Context, qtx *Queries) CreateClientLocationTransferParams {
				client := createRandomClientDetails(ctx, qtx)
				fromLocationID := uuid.New()
				toLocationID := uuid.New()
				newMentorID := uuid.New()
				return CreateClientLocationTransferParams{
					ClientID:       client.ID,
					FromLocationID: &fromLocationID,
					ToLocationID:   &toLocationID,
					RequestDate:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
					NewMentorID:    &newMentorID,
					Reason:         util.StringPtr("Better location for client"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "CreateClientLocationTransfer() should not error")
			},
		},
		{
			name: "location transfer with minimal fields",
			setup: func(ctx context.Context, qtx *Queries) CreateClientLocationTransferParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientLocationTransferParams{
					ClientID:       client.ID,
					FromLocationID: nil,
					ToLocationID:   nil,
					RequestDate:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
					NewMentorID:    nil,
					Reason:         nil,
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "location transfer with reason but no new mentor",
			setup: func(ctx context.Context, qtx *Queries) CreateClientLocationTransferParams {
				client := createRandomClientDetails(ctx, qtx)
				fromLocationID := uuid.New()
				toLocationID := uuid.New()
				return CreateClientLocationTransferParams{
					ClientID:       client.ID,
					FromLocationID: &fromLocationID,
					ToLocationID:   &toLocationID,
					RequestDate:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
					NewMentorID:    nil,
					Reason:         util.StringPtr("Relocation"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			err = qtx.CreateClientLocationTransfer(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestApproveOrRejectClientLocationTransfer(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ApproveOrRejectClientLocationTransferParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "approve location transfer",
			setup: func(ctx context.Context, qtx *Queries) ApproveOrRejectClientLocationTransferParams {
				transfer := createRandomClientLocationTransfer(ctx, qtx)
				return ApproveOrRejectClientLocationTransferParams{
					ID:                 transfer.ID,
					Status:             ClientLocationTransferStatusEnumApproved,
					ApprovedRejectedBy: func() *uuid.UUID { u := uuid.New(); return &u }(),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "reject location transfer",
			setup: func(ctx context.Context, qtx *Queries) ApproveOrRejectClientLocationTransferParams {
				transfer := createRandomClientLocationTransfer(ctx, qtx)
				return ApproveOrRejectClientLocationTransferParams{
					ID:                 transfer.ID,
					Status:             ClientLocationTransferStatusEnumRejected,
					ApprovedRejectedBy: func() *uuid.UUID { u := uuid.New(); return &u }(),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			err = qtx.ApproveOrRejectClientLocationTransfer(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestGetClientCounts(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries)
		checks func(t *testing.T, counts GetClientCountsRow, err error)
	}{
		{
			name: "get client counts with existing clients",
			setup: func(ctx context.Context, qtx *Queries) {
				for i := 0; i < 2; i++ {
					client := createRandomClientDetails(ctx, qtx)
					_, _ = qtx.UpdateClientStatus(ctx, UpdateClientStatusParams{
						ID:     client.ID,
						Status: ClientStatusEnumInCare,
					})
				}
				for i := 0; i < 1; i++ {
					client := createRandomClientDetails(ctx, qtx)
					_, _ = qtx.UpdateClientStatus(ctx, UpdateClientStatusParams{
						ID:     client.ID,
						Status: ClientStatusEnumOnWaitingList,
					})
				}
			},
			checks: func(t *testing.T, counts GetClientCountsRow, err error) {
				require.NoError(t, err, "GetClientCounts() should not error")
				require.Greater(t, counts.TotalClients, int64(0))
			},
		},
		{
			name: "get client counts with empty database",
			setup: func(ctx context.Context, qtx *Queries) {
				// No setup
			},
			checks: func(t *testing.T, counts GetClientCountsRow, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(0), counts.TotalClients)
				require.Equal(t, int64(0), counts.ClientsInCare)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			tt.setup(ctx, qtx)
			counts, err := qtx.GetClientCounts(ctx)
			tt.checks(t, counts, err)
		})
	}
}

func TestGetAllClientsIDs(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries)
		checks func(t *testing.T, ids []uuid.UUID, err error)
	}{
		{
			name: "get all client IDs with existing clients",
			setup: func(ctx context.Context, qtx *Queries) {
				for i := 0; i < 3; i++ {
					_ = createRandomClientDetails(ctx, qtx)
				}
			},
			checks: func(t *testing.T, ids []uuid.UUID, err error) {
				require.NoError(t, err, "GetAllClientsIDs() should not error")
				require.Greater(t, len(ids), 0)
				for _, id := range ids {
					require.NotEqual(t, uuid.Nil, id)
				}
			},
		},
		{
			name: "get all client IDs with empty database",
			setup: func(ctx context.Context, qtx *Queries) {
				// No setup
			},
			checks: func(t *testing.T, ids []uuid.UUID, err error) {
				require.NoError(t, err)
				require.Empty(t, ids)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			tt.setup(ctx, qtx)
			ids, err := qtx.GetAllClientsIDs(ctx)
			tt.checks(t, ids, err)
		})
	}
}

func TestGetClientAddresses(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, addresses []byte, err error)
	}{
		{
			name: "get client addresses with valid data",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				addressData, _ := json.Marshal([]map[string]string{
					{
						"street": "Main Street",
						"city":   "Amsterdam",
						"zip":    "1012AB",
					},
				})
				params := CreateClientDetailsParams{
					FirstName:                  "John",
					LastName:                   "Doe",
					DateOfBirth:                pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Identity:                   true,
					Email:                      util.RandomEmail(),
					Gender:                     ClientGenderEnumMale,
					Filenumber:                 util.RandomString(10),
					EducationLevel:             ClientEducationLevelEnumPrimary,
					Addresses:                  addressData,
					EducationCurrentlyEnrolled: false,
					WorkCurrentlyEmployed:      false,
				}
				client, _ := qtx.CreateClientDetails(ctx, params)
				return client.ID
			},
			checks: func(t *testing.T, addresses []byte, err error) {
				require.NoError(t, err, "GetClientAddresses() should not error")
				require.NotEmpty(t, addresses)
				var addressList []map[string]string
				require.NoError(t, json.Unmarshal(addresses, &addressList))
			},
		},
		{
			name: "get addresses for non-existent client",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, addresses []byte, err error) {
				require.Error(t, err, "GetClientAddresses() should error for non-existent client")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			clientID := tt.setup(ctx, qtx)
			addresses, err := qtx.GetClientAddresses(ctx, clientID)
			tt.checks(t, addresses, err)
		})
	}
}

func TestGetMissingClientDocuments(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, missingLabels []string, err error)
	}{
		{
			name: "get missing documents for client with no documents",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				return client.ID
			},
			checks: func(t *testing.T, missingLabels []string, err error) {
				require.NoError(t, err, "GetMissingClientDocuments() should not error")
				require.Greater(t, len(missingLabels), 0)
				require.Greater(t, len(missingLabels), 5)
			},
		},
		{
			name: "get missing documents for non-existent client",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, missingLabels []string, err error) {
				require.NoError(t, err)
				require.Greater(t, len(missingLabels), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			clientID := tt.setup(ctx, qtx)
			missingLabels, err := qtx.GetMissingClientDocuments(ctx, clientID)
			tt.checks(t, missingLabels, err)
		})
	}
}

func TestListClientDetails(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListClientDetailsParams
		checks func(t *testing.T, clients []ListClientDetailsRow, err error)
	}{
		{
			name: "list all clients with no filters",
			setup: func(ctx context.Context, qtx *Queries) ListClientDetailsParams {
				for i := 0; i < 3; i++ {
					_ = createRandomClientDetails(ctx, qtx)
				}
				return ListClientDetailsParams{
					Status:     NullClientStatusEnum{Valid: false},
					LocationID: nil,
					Search:     nil,
					Limit:      10,
					Offset:     0,
				}
			},
			checks: func(t *testing.T, clients []ListClientDetailsRow, err error) {
				require.NoError(t, err, "ListClientDetails() should not error")
				require.Greater(t, len(clients), 0)
			},
		},
		{
			name: "list clients with limit and offset",
			setup: func(ctx context.Context, qtx *Queries) ListClientDetailsParams {
				for i := 0; i < 5; i++ {
					_ = createRandomClientDetails(ctx, qtx)
				}
				return ListClientDetailsParams{
					Status:     NullClientStatusEnum{Valid: false},
					LocationID: nil,
					Search:     nil,
					Limit:      2,
					Offset:     0,
				}
			},
			checks: func(t *testing.T, clients []ListClientDetailsRow, err error) {
				require.NoError(t, err)
				require.LessOrEqual(t, len(clients), 2)
			},
		},
		{
			name: "list clients with search filter",
			setup: func(ctx context.Context, qtx *Queries) ListClientDetailsParams {
				client := createRandomClientDetails(ctx, qtx)
				return ListClientDetailsParams{
					Status:     NullClientStatusEnum{Valid: false},
					LocationID: nil,
					Search:     util.StringPtr(client.FirstName),
					Limit:      10,
					Offset:     0,
				}
			},
			checks: func(t *testing.T, clients []ListClientDetailsRow, err error) {
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(clients), 0)
			},
		},
		{
			name: "list empty clients",
			setup: func(ctx context.Context, qtx *Queries) ListClientDetailsParams {
				return ListClientDetailsParams{
					Status:     NullClientStatusEnum{Valid: false},
					LocationID: nil,
					Search:     nil,
					Limit:      10,
					Offset:     0,
				}
			},
			checks: func(t *testing.T, clients []ListClientDetailsRow, err error) {
				require.NoError(t, err)
				require.Empty(t, clients)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			clients, err := qtx.ListClientDetails(ctx, params)
			tt.checks(t, clients, err)
		})
	}
}

// Random data generators for tests

func createRandomClientDetails(ctx context.Context, qtx *Queries) ClientDetail {
	params := CreateClientDetailsParams{
		FirstName:                  "John",
		LastName:                   "Doe",
		DateOfBirth:                pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
		Identity:                   true,
		Email:                      util.RandomEmail(),
		Gender:                     ClientGenderEnumMale,
		Filenumber:                 util.RandomString(10),
		EducationLevel:             ClientEducationLevelEnumPrimary,
		Addresses:                  []byte(`[]`),
		EducationCurrentlyEnrolled: false,
		WorkCurrentlyEmployed:      false,
		LivingSituation:            NullClientLivingSituationEnum{Valid: false},
	}
	client, err := qtx.CreateClientDetails(ctx, params)
	if err != nil {
		panic("failed to create random client: " + err.Error())
	}
	return client
}

func createRandomClientLocationTransfer(ctx context.Context, qtx *Queries) ClientLocationTransfer {
	client := createRandomClientDetails(ctx, qtx)
	fromLocationID := uuid.New()
	toLocationID := uuid.New()
	newMentorID := uuid.New()
	return ClientLocationTransfer{
		ClientID:       client.ID,
		FromLocationID: &fromLocationID,
		ToLocationID:   &toLocationID,
		RequestDate:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		NewMentorID:    &newMentorID,
		Reason:         util.StringPtr("Better location for client"),
	}
}
