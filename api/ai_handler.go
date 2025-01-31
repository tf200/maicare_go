package api

import "github.com/gin-gonic/gin"

type CorrectSpellingRequest struct {
	InitialText string `json:"initial_text"`
}

type CorrectSpellingResponse struct {
	CorrectedText string `json:"corrected_text"`
	InitialText   string `json:"initial_text"`
}

func (server *Server) CorrectSpellingApi(ctx *gin.Context) {

}
