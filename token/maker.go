package token

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
	TwoFAToken   TokenType = "2fa_token"
)

func (t TokenType) String() string {
	return string(t)
}

type Maker interface {
	CreateToken(user_id uuid.UUID, employee_id uuid.UUID, duration time.Duration, tokenType TokenType) (string, *Payload, error)
	CreateTokenWithSessionID(user_id uuid.UUID, employee_id uuid.UUID, duration time.Duration, tokenType TokenType, sessionID uuid.UUID) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
