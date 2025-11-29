package care

import "github.com/google/uuid"

// ListCarePlanTopics represents a maturity matrix in the list
type ListCarePlanTopics struct {
	ID        uuid.UUID `json:"id"`
	TopicName string    `json:"topic_name"`
}
