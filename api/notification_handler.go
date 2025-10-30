package api

import (
	"fmt"
	"net/http"

	"maicare_go/service/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	var req notification.ListNotificationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters")))
		return
	}
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}
	res, err := server.businessService.NotificationService.ListNotifications(ctx, &req, payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list notifications")))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Notifications retrieved successfully"))
}

// MarkNotificationAsReadApi handles marking a notification as read
// @Summary Mark Notification as Read
// @Description Marks a notification as read for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} Response[MarkNotificationAsReadResponse]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 403 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /notifications/{id}/read [post]
func (server *Server) MarkNotificationAsReadApi(ctx *gin.Context) {
	notifID := ctx.Param("id")
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}
	parsedNotifID, err := uuid.Parse(notifID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid notification ID format")))
		return
	}
	res, err := server.businessService.NotificationService.MarkNotificationAsRead(ctx, parsedNotifID, payload.UserId)
	if err != nil {
		if err.Error() == "notification does not belong to user" {
			ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("forbidden: notification does not belong to user")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to mark notification as read")))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Notification marked as read successfully"))
}
