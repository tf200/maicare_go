package api

import (
	"database/sql"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

// UpdateClientStatusRequest represents a request to update a client status
type UpdateClientStatusRequest struct {
	Status        string    `json:"status" binding:"required"`
	Reason        string    `json:"reason"`
	IsSchedueled  bool      `json:"schedueled"`
	SchedueledFor time.Time `json:"schedueled_for"`
}

// UpdateClientStatusResponse represents a response to an update client request
type UpdateClientStatusResponse struct {
	ID     int64   `json:"id"`
	Status *string `json:"status"`
}

// UpdateClientStatusApi updates a client
// @Summary Update a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body UpdateClientStatusRequest true "Client status"
// @Success 200 {object} Response[UpdateClientStatusResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/status [put]
func (server *Server) UpdateClientStatusApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateClientStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.IsSchedueled {
		if req.SchedueledFor.Before(time.Now()) {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("schedueled for date must be in the future")))
			return
		}

		schedueledChange, err := server.store.CreateSchedueledClientStatusChange(ctx, db.CreateSchedueledClientStatusChangeParams{
			ClientID:      clientID,
			NewStatus:     &req.Status,
			Reason:        &req.Reason,
			ScheduledDate: pgtype.Date{Time: req.SchedueledFor, Valid: true},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		res := SuccessResponse(UpdateClientStatusResponse{
			ID:     schedueledChange.ClientID,
			Status: schedueledChange.NewStatus,
		}, "Client status update schedueled successfully")
		ctx.JSON(http.StatusOK, res)
		return
	} else {

		tx, err := server.store.ConnPool.Begin(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		defer func() {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
				server.logBusinessEvent(LogLevelError, "UpdateClientStatus", "Failed to rollback db", zap.Error(rollbackErr))
			}
		}()

		qtx := server.store.WithTx(tx)

		oldClient, err := qtx.GetClientDetails(ctx, clientID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		client, err := qtx.UpdateClientStatus(ctx, db.UpdateClientStatusParams{
			ID:     clientID,
			Status: &req.Status,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		_, err = qtx.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
			ClientID:  clientID,
			OldStatus: oldClient.Status,
			NewStatus: req.Status,
			Reason:    &req.Reason,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		res := SuccessResponse(UpdateClientStatusResponse{
			ID:     client.ID,
			Status: client.Status,
		}, "Client status updated successfully")
		ctx.JSON(http.StatusOK, res)
	}

}

// ListStatusHistoryApiResponse represents a response to a list status history request
type ListStatusHistoryApiResponse struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	OldStatus *string   `json:"old_status"`
	NewStatus string    `json:"new_status"`
	ChangedAt time.Time `json:"changed_at"`
	ChangedBy *int64    `json:"changed_by"`
	Reason    *string   `json:"reason"`
}

// ListStatusHistoryApi lists status history of a client
// @Summary List status history of a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[[]ListStatusHistoryApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/status_history [get]
func (server *Server) ListStatusHistoryApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListClientStatusHistoryParams{
		ClientID: clientID,
		Limit:    10,
		Offset:   0,
	}

	statusHistory, err := server.store.ListClientStatusHistory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(statusHistory) == 0 {
		res := SuccessResponse([]string{}, "No status history found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	statusHistoryList := make([]ListStatusHistoryApiResponse, len(statusHistory))
	for i, status := range statusHistory {
		statusHistoryList[i] = ListStatusHistoryApiResponse{
			ID:        status.ID,
			ClientID:  status.ClientID,
			OldStatus: status.OldStatus,
			NewStatus: status.NewStatus,
			ChangedAt: status.ChangedAt.Time,
			ChangedBy: status.ChangedBy,
			Reason:    status.Reason,
		}
	}

	res := SuccessResponse(statusHistoryList, "Status history fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// SetClientProfilePictureRequest represents a request to update a client
type SetClientProfilePictureRequest struct {
	AttachmentID uuid.UUID `json:"attachement_id" binding:"required"`
}

// SetClientProfilePictureResponse represents a response to a set client profile picture request
type SetClientProfilePictureResponse struct {
	ID             int64   `json:"id"`
	ProfilePicture *string `json:"profile_picture"`
}

// SetClientProfilePictureApi sets a client profile picture
// @Summary Set a client profile picture
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body SetClientProfilePictureRequest true "Client profile picture"
// @Success 200 {object} Response[SetClientProfilePictureResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/profile_picture [put]
func (server *Server) SetClientProfilePictureApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req SetClientProfilePictureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.SetClientProfilePictureTxParams{
		ClientID:     clientID,
		AttachmentID: req.AttachmentID,
	}
	client, err := server.store.SetClientProfilePictureTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(SetClientProfilePictureResponse{
		ID:             client.User.ID,
		ProfilePicture: client.User.ProfilePicture,
	}, "Profile picture set successfully")
	ctx.JSON(http.StatusOK, res)

}

// AddClientDocumentApiRequest represents a request to add a document to a client
type AddClientDocumentApiRequest struct {
	AttachmentID uuid.UUID
	Label        string
}

// AddClientDocumentApiResponse represents a response to an add client document request
type AddClientDocumentApiResponse struct {
	ID           int64      `json:"id"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
	ClientID     int64      `json:"client_id"`
	Label        string     `json:"label"`
	Name         string     `json:"name"`
	File         string     `json:"file"`
	Size         int32      `json:"size"`
	IsUsed       bool       `json:"is_used"`
	Tag          *string    `json:"tag"`
	UpdatedAt    time.Time  `json:"updated"`
	CreatedAt    time.Time  `json:"created"`
}

// AddClientDocumentApi adds a document to a client
// @Summary Add a document to a client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body AddClientDocumentApiRequest true "Client document"
// @Success 201 {object} Response[AddClientDocumentApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents [post]
func (server *Server) AddClientDocumentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AddClientDocumentApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddClientDocumentTxParams{
		ClientID:     clientID,
		AttachmentID: req.AttachmentID,
		Label:        req.Label,
	}

	clientDoc, err := server.store.AddClientDocumentTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(AddClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
		ClientID:     clientDoc.ClientDocument.ClientID,
		Label:        clientDoc.ClientDocument.Label,
		Name:         clientDoc.Attachment.Name,
		File:         clientDoc.Attachment.File,
		Size:         clientDoc.Attachment.Size,
		IsUsed:       clientDoc.Attachment.IsUsed,
		Tag:          clientDoc.Attachment.Tag,
		UpdatedAt:    clientDoc.Attachment.Updated.Time,
		CreatedAt:    clientDoc.Attachment.Created.Time,
	}, "Client document added successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientDocumentsApiRequest represents a request to list client documents
type ListClientDocumentsApiRequest struct {
	pagination.Request
}

// ListClientDocumentsApiResponse represents a response to a list client documents request
type ListClientDocumentsApiResponse struct {
	ID             int64      `json:"id"`
	AttachmentUuid *uuid.UUID `json:"attachment_uuid"`
	ClientID       int64      `json:"client_id"`
	Label          string     `json:"label"`
	Uuid           uuid.UUID  `json:"uuid"`
	Name           string     `json:"name"`
	File           *string    `json:"file"`
	Size           int32      `json:"size"`
	IsUsed         bool       `json:"is_used"`
	Tag            *string    `json:"tag"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ListClientDocumentsApi lists documents of a client
// @Summary List documents of a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientDocumentsApiResponse]]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents [get]
func (server *Server) ListClientDocumentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientDocumentsApiRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := req.GetParams()

	clientDocs, err := server.store.ListClientDocuments(ctx, db.ListClientDocumentsParams{
		ClientID: clientID,
		Offset:   params.Offset,
		Limit:    params.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(clientDocs) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientDocumentsApiResponse{}, 0)
		res := SuccessResponse(pag, "No client documents found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	totalCount := clientDocs[0].TotalCount

	clientDocList := make([]ListClientDocumentsApiResponse, len(clientDocs))
	for i, clientDoc := range clientDocs {
		clientDocList[i] = ListClientDocumentsApiResponse{
			ID:             clientDoc.ID,
			AttachmentUuid: clientDoc.AttachmentUuid,
			ClientID:       clientDoc.ClientID,
			Label:          clientDoc.Label,
			Uuid:           clientDoc.Uuid,
			Name:           clientDoc.Name,
			File:           server.generateResponsePresignedURL(&clientDoc.File),
			Size:           clientDoc.Size,
			IsUsed:         clientDoc.IsUsed,
			Tag:            clientDoc.Tag,
			UpdatedAt:      clientDoc.Updated.Time,
			CreatedAt:      clientDoc.Created.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, clientDocList, totalCount)

	res := SuccessResponse(pag, "Client documents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteClientDocumentApiRequest represents a request to delete a client document
type DeleteClientDocumentApiRequest struct {
	AttachmentID uuid.UUID `json:"attachement_id" binding:"required"`
}

// DeleteClientDocumentApiResponse represents a response to a delete client document request
type DeleteClientDocumentApiResponse struct {
	ID           int64      `json:"id"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
}

// DeleteClientDocumentApi deletes a client document
// @Summary Delete a client document
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param document_id path int true "Document ID"
// @Param request body DeleteClientDocumentApiRequest true "Client document"
// @Success 200 {object} Response[DeleteClientDocumentApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/documents/{document_id} [delete]
func (server *Server) DeleteClientDocumentApi(ctx *gin.Context) {
	var req DeleteClientDocumentApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteClientDocumentParams{
		AttachmentID: req.AttachmentID,
	}

	clientDoc, err := server.store.DeleteClientDocumentTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(DeleteClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
	}, "Client document deleted successfully")
	ctx.JSON(http.StatusOK, res)

}

// GetMissingClientDocumentsApiResponse represents a response to a get missing client documents request
type GetMissingClientDocumentsApiResponse struct {
	MissingDocs []string `json:"missing_docs"`
}

// GetMissingClientDocumentsApi gets missing documents of a client
// @Summary Get missing documents of a client
// @Tags clients
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[GetMissingClientDocumentsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/missing_documents [get]
func (server *Server) GetMissingClientDocumentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	missingDocs, err := server.store.GetMissingClientDocuments(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(missingDocs) == 0 {
		res := SuccessResponse(GetMissingClientDocumentsApiResponse{
			MissingDocs: missingDocs,
		}, "No missing client documents found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := SuccessResponse(GetMissingClientDocumentsApiResponse{
		MissingDocs: missingDocs,
	}, "Missing client documents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}
