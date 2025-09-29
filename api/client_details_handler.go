package api

import (
	"fmt"
	_ "maicare_go/pagination" // for swagger
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateClientApi creates a new client
// @Summary Create a new client
// @Tags clients
// @Accept json
// @Produce json
// @Param request body clientp.CreateClientDetailsRequest true "Client details"
// @Success 201 {object} Response[clientp.CreateClientDetailsResponse]
// @Failure 400,404,500 {object} Response[clientp.CreateClientDetailsResponse]
// @Router /clients [post]
func (server *Server) CreateClientApi(ctx *gin.Context) {
	var req clientp.CreateClientDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client, err := server.businessService.ClientService.CreateClientDetails(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(client, "Client created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// ListClientsApi lists clients
// @Summary List clients
// @Tags clients
// @Produce json
// @Param status query string false "Client status"
// @Param location_id query int false "Location ID"
// @Param search query string false "Search query"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListClientsApiResponse]]
// @Failure 400,404,500 {object} Response[clientp.ListClientsApiResponse]
// @Router /clients [get]
func (server *Server) ListClientsApi(ctx *gin.Context) {
	var req clientp.ListClientsApiParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters: %v", err)))
		return
	}

	result, err := server.businessService.ClientService.ListClientDetails(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list clients: %v", err)))
		return
	}

	res := SuccessResponse(result, "Clients fetched successfully")
	ctx.JSON(http.StatusOK, res)

}

// GetClientsCountApi gets the count of clients
// @Summary Get the count of clients
// @Tags clients
// @Produce json
// @Success 200 {object} Response[clientp.GetClientsCountResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/counts [get]
func (server *Server) GetClientsCountApi(ctx *gin.Context) {

	clientCount, err := server.businessService.ClientService.GetClientsCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(clientCount, "Clients count fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientApi gets a client
// @Summary Get a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GetClientApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id} [get]
func (server *Server) GetClientApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID: %v", err)))
		return
	}

	response, err := server.businessService.ClientService.GetClientDetails(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get client details: %v", err)))
		return
	}

	res := SuccessResponse(response, "Client fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetClientAddressesApi gets a client
// @Summary Get a client addresses
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GetClientAddressesApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/addresses [get]
func (server *Server) GetClientAddressesApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	response, err := server.businessService.ClientService.GetClientAddresses(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get client addresses: %v", err)))
		return
	}

	res := SuccessResponse(response, "Client addresses fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateClientApi updates a client
// @Summary Update a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.UpdateClientDetailsRequest true "Client details"
// @Success 200 {object} Response[clientp.UpdateClientDetailsResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id} [put]
func (server *Server) UpdateClientApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID: %v", err)))
		return
	}
	var req clientp.UpdateClientDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedClient, err := server.businessService.ClientService.UpdateClientDetails(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(updatedClient, "Client updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateClientStatusApi updates a client
// @Summary Update a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.UpdateClientStatusRequest true "Client status"
// @Success 200 {object} Response[clientp.UpdateClientStatusResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/status [put]
func (server *Server) UpdateClientStatusApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateClientStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedClient, err := server.businessService.ClientService.UpdateClientStatus(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(updatedClient, "Client status updated successfully")
	ctx.JSON(http.StatusOK, res)

}

// ListStatusHistoryApi lists status history of a client
// @Summary List status history of a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[[]clientp.ListStatusHistoryApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/status_history [get]
func (server *Server) ListStatusHistoryApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	statusHistoryList, err := server.businessService.ClientService.ListStatusHistory(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(statusHistoryList, "Status history fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// SetClientProfilePictureApi sets a client profile picture
// @Summary Set a client profile picture
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.SetClientProfilePictureRequest true "Client profile picture"
// @Success 200 {object} Response[clientp.SetClientProfilePictureResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/profile_picture [put]
func (server *Server) SetClientProfilePictureApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req clientp.SetClientProfilePictureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.ClientService.SetClientProfilePicture(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Profile picture set successfully")
	ctx.JSON(http.StatusOK, res)

}

// AddClientDocumentApi adds a document to a client
// @Summary Add a document to a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.AddClientDocumentApiRequest true "Client document"
// @Success 201 {object} Response[clientp.AddClientDocumentApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents [post]
func (server *Server) AddClientDocumentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.AddClientDocumentApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.AddClientDocument(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client document added successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientDocumentsApi lists documents of a client
// @Summary List documents of a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListClientDocumentsApiResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents [get]
func (server *Server) ListClientDocumentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListClientDocumentsApiRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListClientDocuments(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(pag, "Client documents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteClientDocumentApi deletes a client document
// @Summary Delete a client document
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param document_id path int true "Document ID"
// @Param request body clientp.DeleteClientDocumentApiRequest true "Client document"
// @Success 200 {object} Response[clientp.DeleteClientDocumentApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents/{document_id} [delete]
func (server *Server) DeleteClientDocumentApi(ctx *gin.Context) {
	clientID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req clientp.DeleteClientDocumentApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.DeleteClientDocument(ctx, clientID, req.AttachmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client document deleted successfully")
	ctx.JSON(http.StatusOK, res)

}

// GetMissingClientDocumentsApi gets missing documents of a client
// @Summary Get missing documents of a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GetMissingClientDocumentsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/missing_documents [get]
func (server *Server) GetMissingClientDocumentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.GetMissingClientDocuments(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Missing client documents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}
