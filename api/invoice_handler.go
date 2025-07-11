package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/invoice"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateInvoiceRequest struct {
	ClientID  int64     `json:"client_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func (server *Server) CreateInvoiceHandler(ctx *gin.Context) {
	var req CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	invoiceData := invoice.InvoiceParams{
		ClientID:  req.ClientID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	invoiceResult, err := invoice.GenerateInvoice(server.store, invoiceData, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	invoiceDetailsBytes, err := json.Marshal(invoiceResult.InvoiceDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	invoice, err := server.store.CreateInvoice(ctx.Request.Context(), db.CreateInvoiceParams{
		ClientID:       invoiceResult.ClientID,
		DueDate:        pgtype.Date{Time: time.Now().Add(30 * 24 * time.Hour)},
		TotalAmount:    invoiceResult.TotalAmount,
		InvoiceDetails: invoiceDetailsBytes,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, invoice)

}
