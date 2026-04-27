package repository

import (
	"context"

	"github.com/goccy/go-json"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
)

type MaturityMatrixRepository struct {
	queries *db.Queries
}

func NewMaturityMatrixRepository(queries *db.Queries) *MaturityMatrixRepository {
	return &MaturityMatrixRepository{queries: queries}
}

func (r *MaturityMatrixRepository) ListMaturityMatrix(ctx context.Context) ([]domain.MaturityMatrix, error) {
	rows, err := r.queries.ListCarePlanTopics(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.MaturityMatrix, 0, len(rows))
	for _, row := range rows {
		levels := []domain.MaturityMatrixLevel{}
		if len(row.LevelDescription) > 0 {
			_ = json.Unmarshal(row.LevelDescription, &levels)
		}
		result = append(result, domain.MaturityMatrix{
			ID:                row.ID,
			TopicName:         row.TopicName,
			LevelDescriptions: levels,
		})
	}

	return result, nil
}
