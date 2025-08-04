package grpclient

import (
	context "context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateCarePlan(t *testing.T) {
	req := &PersonalizedCarePlanRequest{
		ClientData: &ClientData{
			Age:              16, // Example age, replace with actual data
			LivingSituation:  "Foster care placement",
			EducationLevel:   "High school",
			DomainName:       "Financiën",
			CurrentLevel:     1, // Example current level, replace with actual data
			LevelDescription: "groeiende complexe schulden",
		},
		DomainDefinitions: map[string]*DomainLevels{
			"Financiën": {
				Levels: map[int32]string{
					1: "groeiende complexe schulden",
					2: "beschikt niet over vrij besteedbaar inkomen of groeiende schulden door spontaan of ongepast uitgeven",
					3: "beschikt over vrij besteedbaar inkomen van ouders zonder verantwoordelijkheid voor noodzakelijke behoeften (zak geld), eventuele schulden zijn stabiel of zijn onder beheer",
					4: "beschikt over vrij besteedbaar inkomen van ouders met enige verantwoordelijkheid voor noodzakelijke behoeften (school geld/lunches), gepast uitgeven, eventuele schulden verminderen",
					5: "beschikt over vrij besteedbaar inkomen (uit klusjes of (bij)baan) met enige verantwoordelijkheid voor noodzakelijke behoeften, aan het eind van de maand is geld over, geen schulden",
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	generatedCarePlan, err := testClient.GenerateCarePlan(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, generatedCarePlan)
	require.NotEmpty(t, generatedCarePlan)
}
