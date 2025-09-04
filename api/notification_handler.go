package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ListNotificationsRequest defines the request structure for listing notifications
type ListNotificationsRequest struct {
	pagination.Request
}

// ListNotificationsResponse defines the response structure for listing notifications
type ListNotificationsResponse struct {
	NotificationID   uuid.UUID `json:"notification_id"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	Data             any       `json:"data"`
	CreatedAT        time.Time `json:"created_at"`
}

// ListNotificationsApi handles the API endpoint for listing notifications
// @Summary List Notifications
// @Description List notifications for the authenticated user
// @Tags Notifications
// @Produce json
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Number of items per page" default(10)
// @Success 200 {object} Response[ListNotificationsResponse] "List of notifications"
// @Failure 400 {object} Response[any] "Invalid request parameters"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /notifications [get]
func (server *Server) ListNotificationsApi(ctx *gin.Context) {
	var req ListNotificationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "ListNotificationsApi", "Failed to bind query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters")))
		return
	}
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListNotificationsApi", "Failed to get auth payload", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}
	params := req.GetParams()
	notifs, err := server.store.ListNotifications(ctx, db.ListNotificationsParams{
		UserID: payload.UserId,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListNotificationsApi", "Failed to list notifications", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list notifications")))
		return
	}

	if len(notifs) == 0 {
		server.logBusinessEvent(LogLevelInfo, "ListNotificationsApi", "No notifications found for user", zap.Int64("user_id", payload.UserId))
		res := SuccessResponse([]ListNotificationsResponse{}, "No notifications found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	var res []ListNotificationsResponse
	for _, notif := range notifs {
		processedData, err := server.notifService.Process(notif.Type, notif.Data)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "ListNotificationsApi", "Failed to process notification data", zap.Error(err), zap.String("notification_type", notif.Type))
			continue
		}

		res = append(res, ListNotificationsResponse{
			NotificationID:   notif.ID,
			NotificationType: notif.Type,
			Message:          notif.Message,
			IsRead:           notif.IsRead,
			Data:             processedData,
			CreatedAT:        notif.CreatedAt.Time,
		})
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Notifications retrieved successfully"))

}

// MarkNotificationAsReadResponse represents the response for marking a notification as read
type MarkNotificationAsReadResponse struct {
	NotificationID   uuid.UUID `json:"notification_id"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	Data             any       `json:"data"`
	CreatedAT        time.Time `json:"created_at"`
}

// MarkNotificationAsReadApi handles marking a notification as read
// @Summary Mark Notification as Read
// @Description Marks a notification as read for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} MarkNotificationAsReadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id}/read [post]
func (server *Server) MarkNotificationAsReadApi(ctx *gin.Context) {
	notifID := ctx.Param("id")
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Failed to get auth payload", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}
	parsedNotifID, err := uuid.Parse(notifID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Invalid notification ID format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid notification ID format")))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Failed to begin transaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to begin transaction")))
		return
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Failed to rollback transaction", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to rollback transaction")))
		}
	}()

	qtx := server.store.WithTx(tx)
	updatedNotif, err := qtx.MarkNotificationAsRead(ctx, parsedNotifID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Failed to mark notification as read", zap.Error(err), zap.String("notification_id", parsedNotifID.String()))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to mark notification as read")))
		return
	}

	// check is the notification belongs to the user
	if updatedNotif.UserID != payload.UserId {
		server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Notification does not belong to user", zap.Int64("user_id", payload.UserId), zap.String("notification_id", parsedNotifID.String()))
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("notification does not belong to user")))
		return
	}
	if err := tx.Commit(ctx); err != nil {
		server.logBusinessEvent(LogLevelError, "SetNotificationAsReadApi", "Failed to commit transaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to commit transaction")))
		return
	}

	processedData, err := server.notifService.Process(updatedNotif.Type, updatedNotif.Data)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListNotificationsApi", "Failed to process notification data", zap.Error(err), zap.String("notification_type", updatedNotif.Type))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process notification data")))
		return
	}

	res := MarkNotificationAsReadResponse{
		NotificationID:   updatedNotif.ID,
		NotificationType: updatedNotif.Type,
		Message:          updatedNotif.Message,
		IsRead:           updatedNotif.IsRead,
		Data:             processedData,
		CreatedAT:        updatedNotif.CreatedAt.Time,
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Notification marked as read successfully"))

}
