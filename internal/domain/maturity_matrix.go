package domain

import (
	"context"

	"github.com/google/uuid"
)

type MaturityMatrixLevel struct {
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MaturityMatrix struct {
	ID                uuid.UUID           `json:"id"`
	TopicName         string              `json:"topic_name"`
	LevelDescriptions []MaturityMatrixLevel `json:"level_descriptions"`
}

type MaturityMatrixService interface {
	ListMaturityMatrix(ctx context.Context) ([]MaturityMatrix, error)
}
