package notification

import (
	"context"
	"database/sql"
	"github.com/goccy/go-json"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *notificationService) ListNotifications(ctx context.Context, req *ListNotificationsRequest, userID uuid.UUID) ([]ListNotificationsResponse, error) {
	params := req.GetParams()
	notifs, err := s.Store.ListNotifications(ctx, db.ListNotificationsParams{
		UserID: userID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListNotifications", "Failed to list notifications", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	response := []ListNotificationsResponse{}
	for _, notif := range notifs {
		var processedData NotificationData
		if err := json.Unmarshal(notif.Data, &processedData); err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListNotifications", "Failed to unmarshal notification data", zap.Error(err), zap.String("notification_id", notif.ID.String()))
			continue
		}

		response = append(response, ListNotificationsResponse{
			NotificationID:   notif.ID,
			NotificationType: string(notif.Type),
			Message:          notif.Message,
			IsRead:           notif.IsRead,
			Data:             processedData,
			CreatedAT:        notif.CreatedAt.Time,
		})
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListNotifications", "Successfully retrieved notifications", zap.Int("count", len(response)), zap.String("user_id", userID.String()))
	return response, nil
}

func (s *notificationService) MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) (*MarkNotificationAsReadResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "MarkNotificationAsRead", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "MarkNotificationAsRead", "Failed to rollback transaction", zap.Error(err))
		}
	}()

	qtx := s.Store.WithTx(tx)
	updatedNotif, err := qtx.MarkNotificationAsRead(ctx, notificationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "MarkNotificationAsRead", "Failed to mark notification as read", zap.Error(err), zap.String("notification_id", notificationID.String()))
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if updatedNotif.UserID != userID {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "MarkNotificationAsRead", "Notification does not belong to user", zap.String("user_id", userID.String()), zap.String("notification_id", notificationID.String()))
		return nil, fmt.Errorf("notification does not belong to user")
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "MarkNotificationAsRead", "Failed to commit transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	response := MarkNotificationAsReadResponse{
		NotificationID:   updatedNotif.ID,
		NotificationType: string(updatedNotif.Type),
		Message:          updatedNotif.Message,
		IsRead:           updatedNotif.IsRead,
		CreatedAT:        updatedNotif.CreatedAt.Time,
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "MarkNotificationAsRead", "Notification marked as read successfully", zap.String("notification_id", notificationID.String()), zap.String("user_id", userID.String()))
	return &response, nil
}
