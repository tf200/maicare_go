package clientp

import (
	"time"

	"github.com/google/uuid"
)

type CreateClientGoalRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	Priority    *string   `json:"priority" binding:"omitempty,oneof=low medium high"`
	TopicID     uuid.UUID `json:"topic_id" binding:"required"`
	SortOrder   *int32    `json:"sort_order"`
}

type CreateClientGoalResponse struct {
	ID                uuid.UUID  `json:"id"`
	ClientID          uuid.UUID  `json:"client_id"`
	Title             string     `json:"title"`
	Description       *string    `json:"description"`
	Priority          string     `json:"priority"`
	Status            string     `json:"status"`
	TopicID           *uuid.UUID `json:"topic_id"`
	TopicNameSnapshot *string    `json:"topic_name_snapshot"`
	Source            string     `json:"source"`
	SortOrder         int32      `json:"sort_order"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type UpdateClientGoalRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	TopicID     *uuid.UUID `json:"topic_id"`
	SortOrder   *int32     `json:"sort_order"`
}

type UpdateClientGoalResponse struct {
	MutationType      string                   `json:"mutation_type"`
	GoalID            uuid.UUID                `json:"goal_id"`
	ReplacementGoalID *uuid.UUID               `json:"replacement_goal_id,omitempty"`
	Goal              CreateClientGoalResponse `json:"goal"`
}
