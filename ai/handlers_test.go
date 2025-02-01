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
