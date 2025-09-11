package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

type SetupTwoFARequest struct {
	UserID int64
}

type SetupTwoFAResult struct {
	QRCodeBase64 string
	Secret       string
}

func (s *authService) SetupTwoFA(req SetupTwoFARequest, ctx context.Context) (*SetupTwoFAResult, error) {
	user, err := s.Store.GetUserByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "SetupTwoFA", "User not found for 2FA setup",
				zap.Int64("userID", req.UserID))
			return nil, ErrUserNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Database error during user retrieval",
			zap.Int64("userID", req.UserID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user.TwoFactorEnabled {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "SetupTwoFA", "2FA already enabled for user",
			zap.Int64("userID", req.UserID))
		return nil, ErrTwoFaAlreadyEnabled
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Maicare",
		AccountName: user.Email,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Error generating 2FA key",
			zap.Int64("userID", req.UserID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate 2FA key")
	}

	secret := key.Secret()

	err = s.Store.CreateTemp2FaSecret(ctx, db.CreateTemp2FaSecretParams{
		ID:                  user.ID,
		TwoFactorSecretTemp: &secret,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Database error saving temp 2FA secret",
			zap.Int64("userID", req.UserID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to save temp 2FA secret")
	}

	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Error generating QR code",
			zap.Int64("userID", req.UserID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate QR code")
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "SetupTwoFA", "2FA setup initiated",
		zap.Int64("userID", req.UserID))

	return &SetupTwoFAResult{
		QRCodeBase64: qrCodeBase64,
		Secret:       secret,
	}, nil

}




