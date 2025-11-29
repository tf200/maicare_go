package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateContract(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateContractParams
		checks func(t *testing.T, contract Contract, err error)
	}{
		{
			name: "successful contract creation",
			setup: func(ctx context.Context, qtx *Queries) CreateContractParams {
				client := createRandomClientDetails(ctx, qtx)
				sender := createRandomSenders(ctx, qtx)
				contractType := createRandomContractType(ctx, qtx)
				return CreateContractParams{
					TypeID:          &contractType.ID,
					Status:          ContractStatusEnumApproved,
					StartDate:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
					EndDate:         pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true},
					ReminderPeriod:  30,
					Vat:             util.Int32Ptr(21),
					Price:           100.0,
					PriceTimeUnit:   PriceTimeUnitEnumHourly,
					Hours:           util.Float64Ptr(40.0),
					HoursType:       HoursTypeEnumWeekly,
					CareName:        util.RandomString(10),
					CareType:        CareTypeEnumAccommodation,
					ClientID:        client.ID,
					SenderID:        &sender.ID,
					AttachmentIds:   []uuid.UUID{uuid.New()},
					FinancingAct:    FinancingActEnumWMO,
					FinancingOption: FinancingOptionEnumPGB,
				}
			},
			checks: func(t *testing.T, contract Contract, err error) {
				require.NoError(t, err, "CreateContract should not return an error")
				require.NotZero(t, contract.ID)
				require.Equal(t, ContractStatusEnumApproved, contract.Status)
			},
		},
		{
			name: "contract creation with minimal fields",
			setup: func(ctx context.Context, qtx *Queries) CreateContractParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateContractParams{
					Status:          ContractStatusEnumApproved,
					StartDate:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
					EndDate:         pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true},
					ReminderPeriod:  30,
					Price:           50.0,
					PriceTimeUnit:   PriceTimeUnitEnumMonthly,
					CareName:        util.RandomString(5),
					CareType:        CareTypeEnumAmbulante,
					ClientID:        client.ID,
					HoursType:       HoursTypeEnumWeekly,
					FinancingAct:    FinancingActEnumWMO,
					FinancingOption: FinancingOptionEnumPGB,
					AttachmentIds:   []uuid.UUID{},
				}
			},
			checks: func(t *testing.T, contract Contract, err error) {
				require.NoError(t, err, "CreateContract should not return an error")
				require.NotZero(t, contract.ID)
				require.Nil(t, contract.TypeID)
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
			contract, err := qtx.CreateContract(ctx, params)
			tt.checks(t, contract, err)
		})
	}
}

func TestCreateContractReminder(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateContractReminderParams
		checks func(t *testing.T, reminder ContractReminder, err error)
	}{
		{
			name: "successful reminder creation",
			setup: func(ctx context.Context, qtx *Queries) CreateContractReminderParams {
				contract := createRandomContract(ctx, qtx)
				return CreateContractReminderParams{
					ContractID:     contract.ID,
					ReminderSentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
				}
			},
			checks: func(t *testing.T, reminder ContractReminder, err error) {
				require.NoError(t, err, "CreateContractReminder should not return an error")
				require.NotZero(t, reminder.ID)
				require.Equal(t, "initial", string(reminder.ReminderType))
			},
		},
		{
			name: "follow-up reminder creation",
			setup: func(ctx context.Context, qtx *Queries) CreateContractReminderParams {
				contract := createRandomContract(ctx, qtx)
				_, _ = qtx.CreateContractReminder(ctx, CreateContractReminderParams{
					ContractID:     contract.ID,
					ReminderSentAt: pgtype.Timestamptz{Time: time.Now().Add(-time.Hour), Valid: true},
				})
				return CreateContractReminderParams{
					ContractID:     contract.ID,
					ReminderSentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
				}
			},
			checks: func(t *testing.T, reminder ContractReminder, err error) {
				require.NoError(t, err, "CreateContractReminder should not return an error")
				require.Equal(t, "follow_up", string(reminder.ReminderType))
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
			reminder, err := qtx.CreateContractReminder(ctx, params)
			tt.checks(t, reminder, err)
		})
	}
}

func TestCreateContractType(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) string
		checks func(t *testing.T, contractType ContractType, err error)
	}{
		{
			name: "successful contract type creation",
			setup: func(ctx context.Context, qtx *Queries) string {
				return util.RandomString(10)
			},
			checks: func(t *testing.T, contractType ContractType, err error) {
				require.NoError(t, err, "CreateContractType should not return an error")
				require.NotZero(t, contractType.ID)
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

			name := tt.setup(ctx, qtx)
			contractType, err := qtx.CreateContractType(ctx, name)
			tt.checks(t, contractType, err)
		})
	}
}

func TestDeleteContractType(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing contract type",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				contractType := createRandomContractType(ctx, qtx)
				return contractType.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteContractType should not error")
			},
		},
		{
			name:  "delete non-existent contract type",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID { return uuid.New() },
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteContractType should not error for non-existent ID")
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

			id := tt.setup(ctx, qtx)
			err = qtx.DeleteContractType(ctx, id)
			tt.checks(t, err)
		})
	}
}

func TestGetBillablePeriodsForContract(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) GetBillablePeriodsForContractParams
		checks func(t *testing.T, periods []GetBillablePeriodsForContractRow, err error)
	}{
		{
			name: "get billable periods for approved contract",
			setup: func(ctx context.Context, qtx *Queries) GetBillablePeriodsForContractParams {
				contract := createRandomContract(ctx, qtx)
				_, _ = qtx.UpdateContractStatus(ctx, UpdateContractStatusParams{
					Status:     ContractStatusEnumApproved,
					ContractID: contract.ID,
				})
				return GetBillablePeriodsForContractParams{
					InvoiceStartDate: pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
					InvoiceEndDate:   pgtype.Timestamptz{Time: time.Now().Add(60 * 24 * time.Hour), Valid: true},
					ContractID:       contract.ID,
				}
			},
			checks: func(t *testing.T, periods []GetBillablePeriodsForContractRow, err error) {
				require.NoError(t, err, "GetBillablePeriodsForContract should not error")
				require.GreaterOrEqual(t, len(periods), 0)
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
			periods, err := qtx.GetBillablePeriodsForContract(ctx, params)
			tt.checks(t, periods, err)
		})
	}
}

func TestGetClientContract(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, row GetClientContractRow, err error)
	}{
		{
			name: "get existing client contract",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				contract := createRandomContract(ctx, qtx)
				return contract.ID
			},
			checks: func(t *testing.T, row GetClientContractRow, err error) {
				require.NoError(t, err, "GetClientContract should not error")
				require.NotZero(t, row.ID)
			},
		},
		{
			name:  "get non-existent client contract",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID { return uuid.New() },
			checks: func(t *testing.T, row GetClientContractRow, err error) {
				require.Error(t, err, "GetClientContract should error for non-existent ID")
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

			id := tt.setup(ctx, qtx)
			row, err := qtx.GetClientContract(ctx, id)
			tt.checks(t, row, err)
		})
	}
}

func TestGetContractAudit(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, audits []GetContractAuditRow, err error)
	}{
		{
			name: "get contract audit",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				contract := createRandomContract(ctx, qtx)
				return contract.ID
			},
			checks: func(t *testing.T, audits []GetContractAuditRow, err error) {
				require.NoError(t, err, "GetContractAudit should not error")
				require.GreaterOrEqual(t, len(audits), 0)
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

			contractID := tt.setup(ctx, qtx)
			audits, err := qtx.GetContractAudit(ctx, contractID)
			tt.checks(t, audits, err)
		})
	}
}

func TestGetSenderContracts(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) *uuid.UUID
		checks func(t *testing.T, contracts []Contract, err error)
	}{
		{
			name: "get sender contracts",
			setup: func(ctx context.Context, qtx *Queries) *uuid.UUID {
				sender := createRandomSenders(ctx, qtx)
				_ = createContractWithSender(ctx, qtx, sender.ID)
				return &sender.ID
			},
			checks: func(t *testing.T, contracts []Contract, err error) {
				require.NoError(t, err, "GetSenderContracts should not error")
				require.GreaterOrEqual(t, len(contracts), 0)
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

			senderID := tt.setup(ctx, qtx)
			contracts, err := qtx.GetSenderContracts(ctx, senderID)
			tt.checks(t, contracts, err)
		})
	}
}

func TestListClientContracts(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListClientContractsParams
		checks func(t *testing.T, contracts []ListClientContractsRow, err error)
	}{
		{
			name: "list client contracts with results",
			setup: func(ctx context.Context, qtx *Queries) ListClientContractsParams {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomContract(ctx, qtx)
				return ListClientContractsParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, contracts []ListClientContractsRow, err error) {
				require.NoError(t, err, "ListClientContracts should not error")
				require.GreaterOrEqual(t, len(contracts), 1)
			},
		},
		{
			name: "list client contracts empty",
			setup: func(ctx context.Context, qtx *Queries) ListClientContractsParams {
				client := createRandomClientDetails(ctx, qtx)
				return ListClientContractsParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, contracts []ListClientContractsRow, err error) {
				require.NoError(t, err, "ListClientContracts should not error")
				require.Empty(t, contracts)
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
			contracts, err := qtx.ListClientContracts(ctx, params)
			tt.checks(t, contracts, err)
		})
	}
}

func TestListContractTypes(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries)
		checks func(t *testing.T, types []ContractType, err error)
	}{
		{
			name: "list contract types",
			setup: func(ctx context.Context, qtx *Queries) {
				_ = createRandomContractType(ctx, qtx)
			},
			checks: func(t *testing.T, types []ContractType, err error) {
				require.NoError(t, err, "ListContractTypes should not error")
				require.GreaterOrEqual(t, len(types), 1)
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

			tt.setup(ctx, qtx)
			types, err := qtx.ListContractTypes(ctx)
			tt.checks(t, types, err)
		})
	}
}

func TestListContracts(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListContractsParams
		checks func(t *testing.T, contracts []ListContractsRow, err error)
	}{
		{
			name: "list contracts with results",
			setup: func(ctx context.Context, qtx *Queries) ListContractsParams {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomContractWithClient(ctx, qtx, client.ID)
				return ListContractsParams{
					Limit:  10,
					Offset: 0,
				}
			},
			checks: func(t *testing.T, contracts []ListContractsRow, err error) {
				require.NoError(t, err, "ListContracts should not error")
				require.GreaterOrEqual(t, len(contracts), 1)
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
			contracts, err := qtx.ListContracts(ctx, params)
			tt.checks(t, contracts, err)
		})
	}
}

func TestListContractsTobeReminded(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries)
		checks func(t *testing.T, contracts []ListContractsTobeRemindedRow, err error)
	}{
		{
			name: "list contracts to be reminded",
			setup: func(ctx context.Context, qtx *Queries) {
				contract := createRandomContract(ctx, qtx)
				_, _ = qtx.UpdateContractStatus(ctx, UpdateContractStatusParams{
					Status:     "approved",
					ContractID: contract.ID,
				})
			},
			checks: func(t *testing.T, contracts []ListContractsTobeRemindedRow, err error) {
				require.NoError(t, err, "ListContractsTobeReminded should not error")
				// May be empty depending on dates
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

			tt.setup(ctx, qtx)
			contracts, err := qtx.ListContractsTobeReminded(ctx)
			tt.checks(t, contracts, err)
		})
	}
}

func TestUpdateContract(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateContractParams
		checks func(t *testing.T, contract Contract, params UpdateContractParams, err error)
	}{
		{
			name: "update existing contract",
			setup: func(ctx context.Context, qtx *Queries) UpdateContractParams {
				contract := createRandomContract(ctx, qtx)
				newCareName := util.RandomString(10)
				return UpdateContractParams{
					ID:       contract.ID,
					CareName: &newCareName,
				}
			},
			checks: func(t *testing.T, contract Contract, params UpdateContractParams, err error) {
				require.NoError(t, err, "UpdateContract should not error")
				require.Equal(t, *params.CareName, contract.CareName)
			},
		},
		{
			name: "update non-existent contract",
			setup: func(ctx context.Context, qtx *Queries) UpdateContractParams {
				return UpdateContractParams{
					ID:       uuid.New(),
					CareName: util.StringPtr(util.RandomString(5)),
				}
			},
			checks: func(t *testing.T, contract Contract, params UpdateContractParams, err error) {
				require.NoError(t, err, "UpdateContract should not error for non-existent ID")
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
			contract, err := qtx.UpdateContract(ctx, params)
			tt.checks(t, contract, params, err)
		})
	}
}

func TestUpdateContractStatus(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateContractStatusParams
		checks func(t *testing.T, contract Contract, err error)
	}{
		{
			name: "update contract status to approved",
			setup: func(ctx context.Context, qtx *Queries) UpdateContractStatusParams {
				contract := createRandomContract(ctx, qtx)
				return UpdateContractStatusParams{
					Status:     ContractStatusEnumApproved,
					ContractID: contract.ID,
				}
			},
			checks: func(t *testing.T, contract Contract, err error) {
				require.NoError(t, err, "UpdateContractStatus should not error")
				require.Equal(t, ContractStatusEnumApproved, contract.Status)
				require.True(t, contract.ApprovedAt.Valid)
			},
		},
		{
			name: "update contract status to draft",
			setup: func(ctx context.Context, qtx *Queries) UpdateContractStatusParams {
				contract := createRandomContract(ctx, qtx)
				return UpdateContractStatusParams{
					Status:     ContractStatusEnumDraft,
					ContractID: contract.ID,
				}
			},
			checks: func(t *testing.T, contract Contract, err error) {
				require.NoError(t, err, "UpdateContractStatus should not error")
				require.Equal(t, "draft", string(contract.Status))
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
			contract, err := qtx.UpdateContractStatus(ctx, params)
			tt.checks(t, contract, err)
		})
	}
}

// Helpers
func createRandomContractType(ctx context.Context, qtx *Queries) ContractType {
	contractType, err := qtx.CreateContractType(ctx, util.RandomString(10))
	if err != nil {
		panic(err)
	}
	return contractType
}

func createRandomContract(ctx context.Context, qtx *Queries) Contract {
	client := createRandomClientDetails(ctx, qtx)

	params := CreateContractParams{
		Status:          ContractStatusEnumApproved,
		StartDate:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true},
		ReminderPeriod:  30,
		Price:           100.0,
		PriceTimeUnit:   PriceTimeUnitEnumMonthly,
		CareName:        util.RandomString(10),
		CareType:        CareTypeEnumAmbulante,
		ClientID:        client.ID,
		FinancingAct:    FinancingActEnumWMO,
		FinancingOption: FinancingOptionEnumPGB,
		HoursType:       HoursTypeEnumWeekly,
		AttachmentIds:   []uuid.UUID{uuid.New()},
	}
	contract, err := qtx.CreateContract(ctx, params)
	if err != nil {
		panic(err)
	}
	return contract
}

func createContractWithSender(ctx context.Context, qtx *Queries, senderID uuid.UUID) Contract {
	client := createRandomClientDetails(ctx, qtx)

	params := CreateContractParams{
		Status:          ContractStatusEnumApproved,
		StartDate:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true},
		ReminderPeriod:  30,
		Price:           100.0,
		PriceTimeUnit:   PriceTimeUnitEnumMonthly,
		CareName:        util.RandomString(10),
		CareType:        CareTypeEnumAmbulante,
		ClientID:        client.ID,
		SenderID:        &senderID,
		FinancingAct:    FinancingActEnumWMO,
		FinancingOption: FinancingOptionEnumPGB,
		HoursType:       HoursTypeEnumWeekly,
		AttachmentIds:   []uuid.UUID{uuid.New()},
	}
	contract, err := qtx.CreateContract(ctx, params)
	if err != nil {
		panic(err)
	}
	return contract
}

func createRandomContractWithClient(ctx context.Context, qtx *Queries, clientID uuid.UUID) Contract {
	params := CreateContractParams{
		Status:          ContractStatusEnumApproved,
		StartDate:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true},
		ReminderPeriod:  30,
		Price:           100.0,
		PriceTimeUnit:   PriceTimeUnitEnumMonthly,
		CareName:        util.RandomString(10),
		CareType:        CareTypeEnumAmbulante,
		ClientID:        clientID,
		FinancingAct:    FinancingActEnumWMO,
		FinancingOption: FinancingOptionEnumPGB,
		HoursType:       HoursTypeEnumWeekly,
		AttachmentIds:   []uuid.UUID{uuid.New()},
	}
	contract, err := qtx.CreateContract(ctx, params)
	if err != nil {
		panic(err)
	}
	return contract
}
