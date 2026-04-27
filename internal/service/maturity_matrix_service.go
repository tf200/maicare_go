package service

import (
	"context"

	"maicare_go/internal/domain"
	"maicare_go/internal/repository"
)

type maturityMatrixService struct {
	repo *repository.MaturityMatrixRepository
}

func NewMaturityMatrixService(repo *repository.MaturityMatrixRepository) domain.MaturityMatrixService {
	return &maturityMatrixService{repo: repo}
}

func (s *maturityMatrixService) ListMaturityMatrix(ctx context.Context) ([]domain.MaturityMatrix, error) {
	return s.repo.ListMaturityMatrix(ctx)
}
