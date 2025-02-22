package db

import (
	"context"
	"maicare_go/util"
	"testing"

	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestListMaturityMatrix(t *testing.T) {

	matrix, err := testQueries.ListMaturityMatrix(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, matrix)
	require.Len(t, matrix, 13)
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

func TestGetClientMaturityMatrixAssessment(t *testing.T) {
	client := createRandomClientDetails(t)
	mma := createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)

	clientMaturityMatrixAssessment, err := testQueries.GetClientMaturityMatrixAssessment(context.Background(), mma.ID)
	require.NoError(t, err)
	require.NotEmpty(t, clientMaturityMatrixAssessment)
	require.Equal(t, mma.ID, clientMaturityMatrixAssessment.ID)
	require.Equal(t, mma.ClientID, clientMaturityMatrixAssessment.ClientID)
	require.Equal(t, mma.MaturityMatrixID, clientMaturityMatrixAssessment.MaturityMatrixID)

}

func createRandomClientGoal(t *testing.T, mmaID int64) ClientGoal {
	arg := CreateClientGoalParams{
		ClientMaturityMatrixAssessmentID: mmaID,
		Description:                      faker.Paragraph(),
		Status:                           "pending",
		TargetLevel:                      int32(rand.Intn(5) + 1),
		StartDate:                        pgtype.Date{Time: util.RandomTIme(), Valid: true},
		TargetDate:                       pgtype.Date{Time: util.RandomTIme(), Valid: true},
	}
	clientGoal, err := testQueries.CreateClientGoal(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientGoal)
	require.Equal(t, arg.ClientMaturityMatrixAssessmentID, clientGoal.ClientMaturityMatrixAssessmentID)
	require.Equal(t, arg.Description, clientGoal.Description)
	require.Equal(t, arg.Status, clientGoal.Status)
	require.Equal(t, arg.TargetLevel, clientGoal.TargetLevel)
	return clientGoal
}

func TestCreateClientGoal(t *testing.T) {
	client := createRandomClientDetails(t)
	mma := createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)
	createRandomClientGoal(t, mma.ID)

}

func TestListClientGoals(t *testing.T) {
	client := createRandomClientDetails(t)
	mma := createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)
	for i := 0; i < 10; i++ {
		createRandomClientGoal(t, mma.ID)
	}

	arg := ListClientGoalsParams{
		ClientMaturityMatrixAssessmentID: mma.ID,
		Limit:                            5,
		Offset:                           5,
	}

	clientGoals, err := testQueries.ListClientGoals(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, clientGoals, 5)

}

func TestGetClientGoal(t *testing.T) {
	client := createRandomClientDetails(t)
	mma := createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)
	goal := createRandomClientGoal(t, mma.ID)

	clientGoal, err := testQueries.GetClientGoal(context.Background(), goal.ID)
	require.NoError(t, err)
	require.NotEmpty(t, clientGoal)
	require.Equal(t, goal.ID, clientGoal.ID)
	require.Equal(t, goal.ClientMaturityMatrixAssessmentID, clientGoal.ClientMaturityMatrixAssessmentID)
	require.Equal(t, goal.Description, clientGoal.Description)
	require.Equal(t, goal.Status, clientGoal.Status)
	require.Equal(t, goal.TargetLevel, clientGoal.TargetLevel)
}

func createRandomGoalObjective(t *testing.T, goalID int64) GoalObjective {
	arg := CreateGoalObjectiveParams{
		GoalID:               goalID,
		ObjectiveDescription: faker.Paragraph(),
		Status:               "pending",
		DueDate:              pgtype.Date{Time: util.RandomTIme(), Valid: true},
	}
	goalObjective, err := testQueries.CreateGoalObjective(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, goalObjective)
	require.Equal(t, arg.GoalID, goalObjective.GoalID)
	require.Equal(t, arg.ObjectiveDescription, goalObjective.ObjectiveDescription)
	require.Equal(t, arg.Status, goalObjective.Status)
	return goalObjective
}

func TestCreateGoalObjective(t *testing.T) {
	client := createRandomClientDetails(t)
	mma := createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)
	goal := createRandomClientGoal(t, mma.ID)
	createRandomGoalObjective(t, goal.ID)
}

func TestListGoalObjectives(t *testing.T) {
	client := createRandomClientDetails(t)
	mma := createRamdomClientMaturityMatrixAssessment(t, client.ID, 1)
	goal := createRandomClientGoal(t, mma.ID)
	for i := 0; i < 10; i++ {
		createRandomGoalObjective(t, goal.ID)
	}

	goalObjectives, err := testQueries.ListGoalObjectives(context.Background(), goal.ID)
	require.NoError(t, err)
	require.Len(t, goalObjectives, 10)
}
