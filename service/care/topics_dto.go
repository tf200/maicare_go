package care

import "github.com/google/uuid"

type LevelDescription struct {
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListCarePlanTopics represents a maturity matrix in the list
type ListCarePlanTopics struct {
	ID                uuid.UUID          `json:"id"`
	TopicName         string             `json:"topic_name"`
	LevelDescriptions []LevelDescription `json:"level_descriptions"`
}
