package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JWTMaker struct {
	accessTokenKey  string
	refreshTokenKey string
	twoFATokenKey   string
}

func NewJWTMaker(accessTokenKey string, refreshTokenKey string, twoFATokenKey string) (Maker, error) {
	if len(accessTokenKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	if len(refreshTokenKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{accessTokenKey, refreshTokenKey, twoFATokenKey}, nil
}

func (maker *JWTMaker) CreateToken(user_id int64, role_id int32, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(user_id, role_id, duration, tokenType)
	if err != nil {
		return "", payload, err
	}

	var secretKey string
	switch tokenType {
	case AccessToken:
		secretKey = maker.accessTokenKey
	case RefreshToken:
		secretKey = maker.refreshTokenKey
	case TwoFAToken:
		secretKey = maker.twoFATokenKey
	default:
		return "", payload, fmt.Errorf("unsupported token type: %v", tokenType)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", payload, fmt.Errorf("failed to create token: %w", err)
	}

	return token, payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		claims, ok := token.Claims.(*Payload)
		if !ok {
			return nil, ErrInvalidToken
		}

		switch claims.TokenType {
		case AccessToken:
			return []byte(maker.accessTokenKey), nil
		case RefreshToken:
			return []byte(maker.refreshTokenKey), nil
		case TwoFAToken:
			return []byte(maker.twoFATokenKey), nil
		default:
			return nil, fmt.Errorf("unknown token type: %v", claims.TokenType)
		}
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
