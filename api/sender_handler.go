package api

import (
	_ "maicare_go/pagination"
	"maicare_go/service/sender"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateSenderApi creates a new sender.
// @Summary Create a new sender
// @Description Create a new sender
// @Tags senders
// @Accept json
// @Produce json
// @Param request body sender.CreateSenderRequest true "Sender data"
// @Success 201 {object} Response[sender.CreateSenderResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders [post]
func (server *Server) CreateSenderApi(ctx *gin.Context) {
	var req sender.CreateSenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.SenderService.CreateSender(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := SuccessResponse(result, "Sender created successfully")

	ctx.JSON(http.StatusCreated, rsp)
}

// ListSendersAPI returns a list of senders.
// @Summary List senders
// @Description List senders
// @Tags senders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search"
// @Param include_archived query bool false "Include archived"
// @Success 200 {object} Response[pagination.Response[sender.ListSendersResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders [get]
func (server *Server) ListSendersAPI(ctx *gin.Context) {
	var req sender.ListSendersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.SenderService.ListSenders(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(response, "Senders retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetSenderAPI returns a sender by ID.
// @Summary Get a sender
// @Description Get a sender
// @Tags senders
// @Produce json
// @Param id path int true "Sender ID"
// @Success 200 {object} Response[sender.GetSenderByIdResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id} [get]
func (server *Server) GetSenderByIdAPI(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.SenderService.GetSenderByID(ctx, senderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Sender retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Update a sender
// @Description Update a sender
// @Tags senders
// @Accept json
// @Produce json
// @Param id path int true "Sender ID"
// @Param request body sender.UpdateSenderRequest true "Sender data"
// @Success 200 {object} Response[sender.UpdateSenderResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id} [put]
func (server *Server) UpdateSenderApi(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req sender.UpdateSenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.SenderService.UpdateSender(ctx, senderID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Sender updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete a sender
// @Description Delete a sender
// @Tags senders
// @Accept json
// @Produce json
// @Param id path int true "Sender ID"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id} [delete]
func (server *Server) DeleteSenderApi(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.SenderService.DeleteSender(ctx, senderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Sender deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateSenderInvoiceTemplateApi handles the creation of a new sender invoice template.
// @Summary Create a new sender invoice template
// @Description Create a new sender invoice template
// @Tags senders
// @Accept json
// @Produce json
// @Param id path int true "Sender ID"
// @Param request body sender.CreateSenderInvoiceTemplateRequest true "Invoice template IDs"
// @Success 201 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /senders/{id}/invoice_template [post]
func (server *Server) CreateSenderInvoiceTemplateApi(ctx *gin.Context) {
	id := ctx.Param("id")
	senderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req sender.CreateSenderInvoiceTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.businessService.SenderService.CreateSenderInvoiceTemplate(ctx, senderID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Sender invoice template created successfully")
	ctx.JSON(http.StatusOK, res)
}
