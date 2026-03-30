package adapters

import (
	"errors"
	"fmt"
	"time"

	"maicare_go/internal/domain"
	pkgjwt "maicare_go/pkg/jwt"

	"github.com/google/uuid"
)

type JWTTokenMakerAdapter struct {
	maker *pkgjwt.Maker
}

func NewJWTTokenMakerAdapter(maker *pkgjwt.Maker) domain.TokenMaker {
	return &JWTTokenMakerAdapter{maker: maker}
}

func (a *JWTTokenMakerAdapter) CreateToken(userID, employeeID uuid.UUID, duration time.Duration, tokenType domain.TokenType) (string, *domain.TokenPayload, error) {
	pkgType, err := toJWTTokenType(tokenType)
	if err != nil {
		return "", nil, err
	}

	token, payload, err := a.maker.CreateToken(userID, employeeID, duration, pkgType)
	if err != nil {
		return "", nil, mapJWTError(err)
	}

	return token, toDomainTokenPayload(payload), nil
}

func (a *JWTTokenMakerAdapter) CreateTokenWithSessionID(userID, employeeID uuid.UUID, duration time.Duration, tokenType domain.TokenType, sessionID uuid.UUID) (string, *domain.TokenPayload, error) {
	pkgType, err := toJWTTokenType(tokenType)
	if err != nil {
		return "", nil, err
	}

	token, payload, err := a.maker.CreateTokenWithSessionID(userID, employeeID, duration, pkgType, sessionID)
	if err != nil {
		return "", nil, mapJWTError(err)
	}

	return token, toDomainTokenPayload(payload), nil
}

func (a *JWTTokenMakerAdapter) VerifyToken(token string) (*domain.TokenPayload, error) {
	payload, err := a.maker.VerifyToken(token)
	if err != nil {
		return nil, mapJWTError(err)
	}

	return toDomainTokenPayload(payload), nil
}

func toJWTTokenType(tokenType domain.TokenType) (pkgjwt.TokenType, error) {
	switch tokenType {
	case domain.AccessTokenType:
		return pkgjwt.AccessToken, nil
	case domain.RefreshTokenType:
		return pkgjwt.RefreshToken, nil
	case domain.TwoFATokenType:
		return pkgjwt.TwoFAToken, nil
	default:
		return "", fmt.Errorf("unsupported token type: %s", tokenType)
	}
}

func toDomainTokenPayload(payload *pkgjwt.Payload) *domain.TokenPayload {
	if payload == nil {
		return nil
	}

	return &domain.TokenPayload{
		ID:         payload.ID,
		SessionID:  payload.SessionID,
		UserID:     payload.UserID,
		EmployeeID: payload.EmployeeID,
		TokenType:  domain.TokenType(payload.TokenType),
		IssuedAt:   payload.IssuedAt,
		ExpiresAt:  payload.ExpiresAt,
	}
}

func mapJWTError(err error) error {
	switch {
	case errors.Is(err, pkgjwt.ErrInvalidToken):
		return domain.ErrInvalidToken
	case errors.Is(err, pkgjwt.ErrExpiredToken):
		return domain.ErrExpiredToken
	default:
		return err
	}
}

var _ domain.TokenMaker = (*JWTTokenMakerAdapter)(nil)
