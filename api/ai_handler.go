package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CorrectSpellingRequest struct {
	InitialText string `json:"initial_text"`
}

type CorrectSpellingResponse struct {
	CorrectedText string `json:"corrected_text"`
	InitialText   string `json:"initial_text"`
}

func (server *Server) SpellingCheckApi(ctx *gin.Context) {
	var request CorrectSpellingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	correctedText, err := server.aiHandler.SpellingCheck(request.InitialText, "mistralai/mistral-small-24b-instruct-2501")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CorrectSpellingResponse{
		CorrectedText: correctedText.CorrectedText,
		InitialText:   request.InitialText,
	}, "Spelling check successful")

	ctx.JSON(http.StatusOK, res)

}
