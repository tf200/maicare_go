package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateClientDetailsRequest struct {
	FirstName             string    `json:"first_name" binding:"required"`
	LastName              string    `json:"last_name" binding:"required"`
	Email                 string    `json:"email" binding:"required,email"`
	Organisation          string    `json:"organisation" binding:"required"`
	Location              *int64    `json:"location" binding:"required"`
	LegalMeasure          string    `json:"legal_measure"`
	Birthplace            string    `json:"birthplace" binding:"required"`
	Departement           string    `json:"departement" binding:"required"`
	Gender                string    `json:"gender" binding:"required"`
	Filenumber            int32     `json:"filenumber" binding:"required"`
	DateOfBirth           string    `json:"date_of_birth" binding:"required"`
	PhoneNumber           *string   `json:"phone_number" binding:"required"`
	Infix                 *string   `json:"infix"`
	Bsn                   *string   `json:"bsn"`
	Source                string    `json:"source" binding:"required"`
	Addresses             []Address `json:"addresses"`
	IdentityAttachmentIds []string  `json:"identity_attachment_ids"`
}

type Address struct {
	BelongsTo   *string `json:"belongs_to"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	ZipCode     *string `json:"zip_code"`
	PhoneNumber *string `json:"phone_number"`
}


func CreateClientApi(ctx *gin.Context) {
	var CreateClientRequest CreateClientDetailsRequest
	if err := ctx.ShouldBindJSON(&CreateClientRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	

}	
		
