package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CorrectSpellingRequest is the request format for the spelling check API
type CorrectSpellingRequest struct {
	InitialText string `json:"initial_text"`
}

// CorrectSpellingResponse is the response format for the spelling check API
type CorrectSpellingResponse struct {
	CorrectedText string `json:"corrected_text"`
	InitialText   string `json:"initial_text"`
}

// SpellingCheckApi is the handler for the spelling check API
// @Summary Check spelling of a text
// @Description Check spelling of a text
// @Tags ai
// @Accept json
// @Produce json
// @Param request body CorrectSpellingRequest true "Request body"
// @Success 200 {object} Response[CorrectSpellingResponse]
// @Router /ai/spelling_check [post]
func (server *Server) SpellingCheckApi(ctx *gin.Context) {
	var request CorrectSpellingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		server.logBusinessEvent(LogLevelWarn, "SpellingCheckApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Printf("Received request: %+v", request)

	correctedText, err := server.aiHandler.SpellingCheck(request.InitialText, "anthropic/claude-3.5-haiku-20241022:beta")
	if err != nil {
		server.logBusinessEvent(LogLevelError, "SpellingCheckApi", "Failed to correct spelling", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CorrectSpellingResponse{
		CorrectedText: correctedText.CorrectedText,
		InitialText:   request.InitialText,
	}, "Spelling check successful")

	server.logBusinessEvent(LogLevelInfo, "SpellingCheckApi", "Spelling check completed successfully", zap.String("corrected_text", correctedText.CorrectedText))

	ctx.JSON(http.StatusOK, res)
}
