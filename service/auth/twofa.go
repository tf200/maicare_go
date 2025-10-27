package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

func (s *authService) SetupTwoFA(userID uuid.UUID, ctx context.Context) (*Setup2FAResponse, error) {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "SetupTwoFA", "User not found for 2FA setup",
				zap.String("userID", userID.String()))
			return nil, ErrUserNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Database error during user retrieval",
			zap.String("userID", userID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user.TwoFactorEnabled {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "SetupTwoFA", "2FA already enabled for user",
			zap.String("userID", userID.String()))
		return nil, ErrTwoFaAlreadyEnabled
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Maicare",
		AccountName: user.Email,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Error generating 2FA key",
			zap.String("userID", userID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate 2FA key")
	}

	secret := key.Secret()

	err = s.Store.CreateTemp2FaSecret(ctx, db.CreateTemp2FaSecretParams{
		ID:                  user.ID,
		TwoFactorSecretTemp: &secret,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Database error saving temp 2FA secret",
			zap.String("userID", userID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to save temp 2FA secret")
	}

	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetupTwoFA", "Error generating QR code",
			zap.String("userID", userID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate QR code")
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "SetupTwoFA", "2FA setup initiated",
		zap.String("userID", userID.String()))

	return &Setup2FAResponse{
		QrCode: qrCodeBase64,
		Secret: secret,
	}, nil

}

func (s *authService) EnableTwoFA(req Enable2FARequest, userID uuid.UUID, ctx context.Context) (*Enable2FAResponse, error) {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "EnableTwoFA", "User not found for 2FA enable",
				zap.String("userID", userID.String()))
			return nil, ErrUserNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "EnableTwoFA", "Database error during user retrieval",
			zap.String("userID", userID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user.TwoFactorEnabled {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "EnableTwoFA", "2FA already enabled for user",
			zap.String("userID", userID.String()))
		return nil, ErrTwoFaAlreadyEnabled
	}
	if user.TwoFactorSecretTemp == nil || *user.TwoFactorSecretTemp == "" {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "EnableTwoFA", "No temp 2FA secret found for user",
			zap.String("userID", userID.String()))
		return nil, fmt.Errorf("no temp 2FA secret found")
	}
	valid := totp.Validate(req.ValidationCode, *user.TwoFactorSecretTemp)
	if !valid {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "EnableTwoFA", "Invalid 2FA validation code",
			zap.String("userID", userID.String()))
		return nil, ErrInvalidTwoFACode
	}

	recoveryCodes := util.GenerateRecoveryCodes(10)

	hashedRecoveryCodes := make([]string, len(recoveryCodes))
	for i, code := range recoveryCodes {
		hashedCode, err := util.HashPassword(code)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "EnableTwoFA", "Error hashing recovery code",
				zap.String("userID", userID.String()), zap.String("error", err.Error()))
			return nil, fmt.Errorf("failed to hash recovery code: %v", err)
		}
		hashedRecoveryCodes[i] = hashedCode
	}

	err = s.Store.Enable2Fa(ctx, db.Enable2FaParams{
		ID:              user.ID,
		TwoFactorSecret: user.TwoFactorSecretTemp,
		RecoveryCodes:   hashedRecoveryCodes,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "EnableTwoFA", "Database error enabling 2FA",
			zap.String("userID", userID.String()), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to enable 2FA: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "EnableTwoFA", "2FA enabled successfully",
		zap.String("userID", userID.String()))
	return &Enable2FAResponse{
		RecoveryCodes: recoveryCodes,
	}, nil
}
