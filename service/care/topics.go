package care

import (
	"context"
	"fmt"

	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *carePlanService) ListCarePlanTopics(ctx context.Context) ([]ListCarePlanTopics, error) {
	topics, err := s.Store.ListMaturityMatrix(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListCarePlanTopics", "Failed to list care plan topics", zap.Error(err))
		return nil, fmt.Errorf("failed to list care plan topics")
	}

	reponse := []ListCarePlanTopics{}
	for _, topic := range topics {
		reponse = append(reponse, ListCarePlanTopics{
			ID:        topic.ID,
			TopicName: topic.TopicName,
		})
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListCarePlanTopics", "Successfully listed care plan topics")
	return reponse, nil
}
