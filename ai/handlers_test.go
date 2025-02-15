package ai

import (
	"maicare_go/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpellingCheck(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	a := NewAiHandler(config.OpenRouterAPIKey)
	text := "The qwick brown fox jumpped over the laazy dog. It was a beutiful day, but sudenly, the wheather changed dramaticaly. Peopel ran for sheltr as the rain began pooring down. In the distanse, a lound thunder clap made everyonne jump. This storm came out of no where! exclaimed a passerby."
	model := "mistralai/mistral-small-24b-instruct-2501"

	response, err := a.SpellingCheck(text, model)
	require.NoError(t, err)
	require.NotNil(t, response)

	// Print response to console
	t.Logf("Spell check response: %s", response)
}

func TestGenerateAutoReports(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	a := NewAiHandler(config.OpenRouterAPIKey)
	text := `Medical Progression Report
			Patient Name: Adam Benali
			Date: April 15, 2024
			Age: 4 years, 3 months
			Report Summary:
			Adam has shown steady growth, now weighing 17 kg and measuring 105 cm in height. 
			His motor skills are developing well, with improved balance and coordination. 
			No major health issues were reported, aside from a minor cold two weeks prior, which resolved without complications. 
			Vaccinations are up to date. 
			Parents note increased curiosity and improved social interactions. 
			Recommendation: 
			Maintain a balanced diet and regular outdoor play. 
			Next check-up in six months.
			Medical Progression Report
			Patient Name: Adam Benali
			Date: January 10, 2024
			Age: 4 years
			Report Summary:
			Adam was brought in for his routine pediatric check-up. His weight is 16 kg, and height is 102 cm, both within the normal range for his age. Parents report good appetite and normal sleeping patterns. His speech development has improved significantly, with clear articulation of sentences. Mild seasonal allergies noted, and an antihistamine was prescribed. No concerns regarding motor skills or cognitive abilities. Follow-up in three months.
			Medical Progression Report
			Patient Name: Adam Benali
			Date: April 15, 2024
			Age: 4 years, 3 months
			Report Summary:
			Adam has shown steady growth, now weighing 17 kg and measuring 105 cm in height. His motor skills are developing well, with improved balance and coordination. No major health issues were reported, aside from a minor cold two weeks prior, which resolved without complications. Vaccinations are up to date. Parents note increased curiosity and improved social interactions. Recommendation: Maintain a balanced diet and regular outdoor play. Next check-up in six months.
			`
	model := "mistralai/mistral-small-24b-instruct-2501"

	response, err := a.GenerateAutoReports(text, model)
	require.NoError(t, err)
	require.NotNil(t, response)

	// Print response to console
	t.Logf("Spell check response: %s", response)
}
