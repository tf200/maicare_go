package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomContractType(t *testing.T) ContractType {
	// Create a random contract type

	contractType, err := testQueries.CreateContractType(context.Background(), "Test Contract Type")
	require.NoError(t, err)
	require.NotEmpty(t, contractType)
	require.NotEmpty(t, contractType.ID)
	require.Equal(t, "Test Contract Type", contractType.Name)
	return contractType
}

func TestCreateContractType(t *testing.T) {
	createRandomContractType(t)

}

func TestListContractType(t *testing.T) {
	// Create 10 random contract types
	for i := 0; i < 10; i++ {
		createRandomContractType(t)
	}

	contractTypes, err := testQueries.ListContractTypes(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, contractTypes)
	require.Len(t, contractTypes, 10)
}

func createRandomContract(t *testing.T, clientID int64, senderID *int64) Contract {
	priceFrequency := []string{"minute", "hourly", "daily", "weekly", "monthly"}
	hoursType := []string{"weekly", "all_period"}
	careType := []string{"ambulante", "accommodation"}
	financingAct := []string{"WMO", "ZVW", "WLZ", "JW", "WPG"}
	financingOption := []string{"ZIN", "PGB"}

	contractType := createRandomContractType(t)
	attachment := createRandomAttachmentFile(t)

	arg := CreateContractParams{
		TypeID:          &contractType.ID,
		StartDate:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Now().AddDate(0, 2, 0), Valid: true},
		ReminderPeriod:  10,
		Tax:             util.Int32Ptr(15),
		Price:           5.58,
		PriceFrequency:  util.RandomEnum(priceFrequency),
		Hours:           util.Int32Ptr(100),
		HoursType:       util.RandomEnum(hoursType),
		CareName:        "Test Care",
		CareType:        util.RandomEnum(careType),
		ClientID:        clientID,
		SenderID:        senderID,
		FinancingAct:    util.RandomEnum(financingAct),
		FinancingOption: util.RandomEnum(financingOption),
		AttachmentIds:   []uuid.UUID{attachment.Uuid},
		Status:          "approved",
	}

	contract, err := testQueries.CreateContract(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contract)
	require.Equal(t, arg.TypeID, contract.TypeID)
	require.Equal(t, arg.ReminderPeriod, contract.ReminderPeriod)
	require.Equal(t, arg.Tax, contract.Tax)
	require.Equal(t, arg.Price, contract.Price)
	require.Equal(t, arg.PriceFrequency, contract.PriceFrequency)
	require.Equal(t, arg.Hours, contract.Hours)
	require.Equal(t, arg.HoursType, contract.HoursType)
	require.Equal(t, arg.CareName, contract.CareName)
	require.Equal(t, arg.CareType, contract.CareType)
	require.Equal(t, arg.ClientID, contract.ClientID)
	require.Equal(t, arg.SenderID, contract.SenderID)
	require.Equal(t, arg.FinancingAct, contract.FinancingAct)
	require.Equal(t, arg.FinancingOption, contract.FinancingOption)
	require.NotZero(t, contract.ID)
	return contract
}

func TestCreateContract(t *testing.T) {
	client := createRandomClientDetails(t)

	createRandomContract(t, client.ID, client.SenderID)
}

func TestUpdateContract(t *testing.T) {
	client := createRandomClientDetails(t)
	contract := createRandomContract(t, client.ID, client.SenderID)

	arg := UpdateContractParams{
		ID:             contract.ID,
		StartDate:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		EndDate:        pgtype.Timestamptz{Time: time.Now().AddDate(0, 2, 0), Valid: true},
		ReminderPeriod: util.Int32Ptr(10),
		Tax:            util.Int32Ptr(15),
		PriceFrequency: util.StringPtr("monthly"),
		Hours:          util.Int32Ptr(100),
		HoursType:      util.StringPtr("all_period"),
		CareName:       util.StringPtr("Test Care"),
		CareType:       util.StringPtr("accommodation"),
		SenderID:       client.SenderID,
	}

	contract2, err := testQueries.UpdateContract(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contract2)
	require.Equal(t, arg.ID, contract2.ID)
	require.Equal(t, arg.ReminderPeriod, contract2.ReminderPeriod)
	require.Equal(t, arg.Tax, contract2.Tax)
	require.Equal(t, arg.Price, contract2.Price)
	require.Equal(t, arg.PriceFrequency, contract2.PriceFrequency)
	require.Equal(t, arg.Hours, contract2.Hours)
	require.Equal(t, arg.HoursType, contract2.HoursType)
	require.Equal(t, arg.CareName, contract2.CareName)
	require.Equal(t, arg.CareType, contract2.CareType)
	require.Equal(t, arg.SenderID, contract2.SenderID)
	require.Equal(t, arg.FinancingAct, contract2.FinancingAct)
	require.Equal(t, arg.FinancingOption, contract2.FinancingOption)
}

func TestGetClientContract(t *testing.T) {
	client := createRandomClientDetails(t)
	contract := createRandomContract(t, client.ID, client.SenderID)
	contract2, err := testQueries.GetClientContract(context.Background(), contract.ID)
	require.NoError(t, err)
	require.NotEmpty(t, contract2)
	require.Equal(t, contract.ID, contract2.ID)
	require.Equal(t, contract.TypeID, contract2.TypeID)
	require.Equal(t, contract.ReminderPeriod, contract2.ReminderPeriod)
	require.Equal(t, contract.Tax, contract2.Tax)
	require.Equal(t, contract.Price, contract2.Price)
	require.Equal(t, contract.PriceFrequency, contract2.PriceFrequency)
	require.Equal(t, contract.Hours, contract2.Hours)
	require.Equal(t, contract.HoursType, contract2.HoursType)
	require.Equal(t, contract.CareName, contract2.CareName)
	require.Equal(t, contract.CareType, contract2.CareType)
	require.Equal(t, contract.ClientID, contract2.ClientID)
	require.Equal(t, contract.SenderID, contract2.SenderID)
	require.Equal(t, contract.FinancingAct, contract2.FinancingAct)
	require.Equal(t, contract.FinancingOption, contract2.FinancingOption)
	require.Equal(t, contract.AttachmentIds, contract2.AttachmentIds)
	require.Equal(t, contract.DepartureReason, contract2.DepartureReason)
	require.Equal(t, contract.DepartureReport, contract2.DepartureReport)
}

func TestListClientContracts(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 10; i++ {
		createRandomContract(t, client.ID, client.SenderID)
	}

	contracts, err := testQueries.ListClientContracts(context.Background(), ListClientContractsParams{
		ClientID: client.ID,
		Limit:    5,
		Offset:   0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, contracts)
	require.Len(t, contracts, 5)

}

func TestListContracts(t *testing.T) {
	// Create 10 random contracts
	for i := 0; i < 10; i++ {
		client := createRandomClientDetails(t)
		createRandomContract(t, client.ID, client.SenderID)
	}

	contracts, err := testQueries.ListContracts(context.Background(), ListContractsParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, contracts)
	require.Len(t, contracts, 5)

}
