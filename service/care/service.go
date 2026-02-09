package care

import (
	"context"

	"maicare_go/service/deps"
)

type CarePlanService interface {
	ListCarePlanTopics(ctx context.Context) ([]ListCarePlanTopics, error)
}

type carePlanService struct {
	*deps.ServiceDependencies
}

func NewCarePlanService(deps *deps.ServiceDependencies) CarePlanService {
	return &carePlanService{
		ServiceDependencies: deps,
	}
}
