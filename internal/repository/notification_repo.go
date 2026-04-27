package repository

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
)

type NotificationRepository struct {
	queries *db.Queries
	store   *db.Store
}

func NewNotificationRepository(queries *db.Queries, store *db.Store) domain.NotificationRepository {
	return &NotificationRepository{queries: queries, store: store}
}

func (r *NotificationRepository) ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Notification, error) {
	rows, err := r.queries.ListNotifications(ctx, db.ListNotificationsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]domain.Notification, 0, len(rows))
	for _, row := range rows {
		var data domain.NotificationData
		if len(row.Data) > 0 {
			_ = json.Unmarshal(row.Data, &data)
		}
		result = append(result, domain.Notification{
			ID:        row.ID,
			UserID:    row.UserID,
			Type:      string(row.Type),
			Message:   row.Message,
			IsRead:    row.IsRead,
			Data:      data,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return result, nil
}

func (r *NotificationRepository) MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) (*domain.Notification, error) {
	tx, err := r.store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := r.store.WithTx(tx)
	row, err := qtx.MarkNotificationAsRead(ctx, notificationID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var data domain.NotificationData
	if len(row.Data) > 0 {
		_ = json.Unmarshal(row.Data, &data)
	}

	return &domain.Notification{
		ID:        row.ID,
		UserID:    row.UserID,
		Type:      string(row.Type),
		Message:   row.Message,
		IsRead:    row.IsRead,
		Data:      data,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *NotificationRepository) CountNotificationsByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.store.ConnPool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE user_id = $1", userID).Scan(&count)
	return count, err
}
