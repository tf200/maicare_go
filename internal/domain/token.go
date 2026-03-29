package domain

import pkgjwt "maicare_go/pkg/jwt"

type TokenVerifier interface {
	VerifyToken(token string) (*pkgjwt.Payload, error)
}
