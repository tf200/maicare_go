package token

import "time"

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

func (t TokenType) String() string {
	return string(t)
}

type Maker interface {
	CreateToken(user_id int64, role_id int32, duration time.Duration, tokenType TokenType) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
