package care

import (
	"context"
	"github.com/goccy/go-json"
	"fmt"

	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *carePlanService) ListCarePlanTopics(ctx context.Context) ([]ListCarePlanTopics, error) {
	topics, err := s.Store.ListCarePlanTopics(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListCarePlanTopics", "Failed to list care plan topics", zap.Error(err))
		return nil, fmt.Errorf("failed to list care plan topics")
	}

	reponse := []ListCarePlanTopics{}
	for _, topic := range topics {
		levelDescriptions := []LevelDescription{}
		err = json.Unmarshal(topic.LevelDescription, &levelDescriptions)
		if err := json.Unmarshal(topic.LevelDescription, &levelDescriptions); err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListCarePlanTopics",
				"Failed to unmarshal level descriptions",
				zap.Error(err),
				zap.ByteString("raw", topic.LevelDescription),
			)
		}
		reponse = append(reponse, ListCarePlanTopics{
			ID:                topic.ID,
			TopicName:         topic.TopicName,
			LevelDescriptions: levelDescriptions,
		})
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListCarePlanTopics", "Successfully listed care plan topics")
	return reponse, nil
}
