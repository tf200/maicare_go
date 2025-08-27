package db

import (
	"context"
	"maicare_go/util"
	"strconv"
	"testing"

	"math/rand"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestListMaturityMatrix(t *testing.T) {

	matrix, err := testQueries.ListMaturityMatrix(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, matrix)
	require.Len(t, matrix, 13)
}

func TestGetLevelDescription(t *testing.T) {
	arg := GetLevelDescriptionParams{
		ID:    1,
		Level: strconv.Itoa(1),
	}
	levelDescription, err := testQueries.GetLevelDescription(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, levelDescription)
}

func createRamdomClientMaturityMatrixAssessment(t *testing.T, clientID int64, maturityMatrixID int64) CreateClientMaturityMatrixAssessmentRow {
	startDate := util.RandomTIme()
	endDate := startDate.AddDate(0, 0, 7)
	arg := CreateClientMaturityMatrixAssessmentParams{
		ClientID:         clientID,
		MaturityMatrixID: maturityMatrixID,
		StartDate:        pgtype.Date{Time: startDate, Valid: true},
		EndDate:          pgtype.Date{Time: endDate, Valid: true},
		InitialLevel:     int32(rand.Intn(5) + 1),
		CurrentLevel:     int32(rand.Intn(5) + 1),
	}

	clientMaturityMatrixAssessment, err := testQueries.CreateClientMaturityMatrixAssessment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientMaturityMatrixAssessment)
	require.Equal(t, arg.ClientID, clientMaturityMatrixAssessment.ClientID)
	require.Equal(t, arg.MaturityMatrixID, clientMaturityMatrixAssessment.MaturityMatrixID)
	require.NotEmpty(t, clientMaturityMatrixAssessment.TopicName)
	return clientMaturityMatrixAssessment
}

func TestCreateClientMaturityMatrixAssessment(t *testing.T) {
	client := createRandomClientDetails(t)
	createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)
}

func TestListClientMaturityMatrixAssessments(t *testing.T) {
	client := createRandomClientDetails(t)
	/*
		we are using i instead of random value because of the unique constraint
		each client can have only one assessment for each maturity matrix
		this way we get a sure unique value between 1 and 10
	*/
	var i int64
	for i = 0; i < 10; i++ {
		createRamdomClientMaturityMatrixAssessment(t, client.ID, i+1)
		arg := ListClientMaturityMatrixAssessmentsParams{
			ClientID: client.ID,
			Limit:    5,
			Offset:   5,
		}
		clientMaturityMatrixAssessments, err := testQueries.ListClientMaturityMatrixAssessments(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, clientMaturityMatrixAssessments, 5)
	}
}
