package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID         uuid.UUID
	UserId     int64     `json:"user_id"`
	EmployeeID int64     `json:"employee_id"`
	TokenType  TokenType `json:"token_type"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// You can adjust this duration

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewPayload(user_id int64, employee_id int64, duration time.Duration, tokenType TokenType) (*Payload, error) {

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	payload := &Payload{
		ID:         tokenID,
		EmployeeID: employee_id,
		UserId:     user_id,
		TokenType:  tokenType,
		IssuedAt:   now,
		ExpiresAt:  now.Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
