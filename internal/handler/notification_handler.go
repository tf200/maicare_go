package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"
	"maicare_go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service domain.NotificationService
}

func NewNotificationHandler(service domain.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) ListNotifications(ctx *gin.Context) {
	payload, ok := middleware.AuthPayloadFromContext(ctx.Request.Context())
	if !ok || payload == nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized access", ""))
		return
	}

	var req listNotificationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid query parameters", ""))
		return
	}

	params := req.Params()
	result, err := h.service.ListNotifications(ctx.Request.Context(), payload.UserID, params.Limit, params.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list notifications", ""))
		return
	}

	items := make([]notificationResponse, len(result))
	for i, n := range result {
		items[i] = toNotificationResponse(n)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(items, "Notifications retrieved successfully"))
}

func (h *NotificationHandler) MarkNotificationAsRead(ctx *gin.Context) {
	payload, ok := middleware.AuthPayloadFromContext(ctx.Request.Context())
	if !ok || payload == nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized access", ""))
		return
	}

	notifID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid notification ID format", ""))
		return
	}

	notif, err := h.service.MarkNotificationAsRead(ctx.Request.Context(), notifID, payload.UserID)
	if err != nil {
		if err.Error() == "notification does not belong to user" {
			ctx.JSON(http.StatusForbidden, httpapi.Fail("forbidden: notification does not belong to user", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to mark notification as read", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toNotificationResponse(*notif), "Notification marked as read successfully"))
}
