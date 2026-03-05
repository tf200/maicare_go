package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

func (s *authService) SetupTwoFA(req Setup2FARequest, userID uuid.UUID, ctx context.Context) (*Setup2FAResponse, error) {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogError(ctx, "SetupTwoFA", "User not found for 2FA setup", nil,
				zap.String("user_id", userID.String()))
			return nil, ErrUserNotFound
		}
		s.Logger.LogError(ctx, "SetupTwoFA", "Database error during user retrieval", err,
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user.TwoFactorEnabled {
		s.Logger.LogWarn(ctx, "SetupTwoFA", "2FA already enabled for user",
			zap.String("user_id", userID.String()))
		return nil, ErrTwoFaAlreadyEnabled
	}

	err = util.CheckPassword(req.CurrentPassword, user.Password)
	if err != nil {
		s.Logger.LogWarn(ctx, "SetupTwoFA", "Invalid password for 2FA setup",
			zap.String("user_id", userID.String()))
		return nil, ErrInvalidCredentials
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Maicare",
		AccountName: user.Email,
	})
	if err != nil {
		s.Logger.LogError(ctx, "SetupTwoFA", "Error generating 2FA key", err,
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to generate 2FA key")
	}

	secret := key.Secret()

	rows, err := s.Store.CreateTemp2FaSecret(ctx, db.CreateTemp2FaSecretParams{
		ID:                  user.ID,
		TwoFactorSecretTemp: &secret,
	})
	if err != nil || rows != 1 {
		s.Logger.LogError(ctx, "SetupTwoFA", "Database error saving temp 2FA secret", err,
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to save temp 2FA secret")
	}

	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		s.Logger.LogError(ctx, "SetupTwoFA", "Error generating QR code", err,
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to generate QR code")
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	s.Logger.LogInfo(ctx, "SetupTwoFA", "2FA setup initiated",
		zap.String("user_id", userID.String()))

	return &Setup2FAResponse{
		QrCode: qrCodeBase64,
		Secret: secret,
	}, nil
}

func (s *authService) EnableTwoFA(req Enable2FARequest, userID uuid.UUID, ctx context.Context) (*Enable2FAResponse, error) {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogWarn(ctx, "EnableTwoFA", "User not found for 2FA enable",
				zap.String("user_id", userID.String()))
			return nil, ErrUserNotFound
		}
		s.Logger.LogError(ctx, "EnableTwoFA", "Database error during user retrieval", err,
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user.TwoFactorEnabled {
		s.Logger.LogWarn(ctx, "EnableTwoFA", "2FA already enabled for user",
			zap.String("user_id", userID.String()))
		return nil, ErrTwoFaAlreadyEnabled
	}
	if user.TwoFactorSecretTemp == nil || *user.TwoFactorSecretTemp == "" {
		s.Logger.LogWarn(ctx, "EnableTwoFA", "No temp 2FA secret found for user",
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("no temp 2FA secret found")
	}
	valid := totp.Validate(req.ValidationCode, *user.TwoFactorSecretTemp)
	if !valid {
		s.Logger.LogWarn(ctx, "EnableTwoFA", "Invalid 2FA validation code",
			zap.String("user_id", userID.String()))
		return nil, ErrInvalidTwoFACode
	}

	recoveryCodes := util.GenerateRecoveryCodes(10)

	hashedRecoveryCodes := make([]string, len(recoveryCodes))
	for i, code := range recoveryCodes {
		hashedCode, err := util.HashPassword(code)
		if err != nil {
			s.Logger.LogError(ctx, "EnableTwoFA", "Error hashing recovery code", err,
				zap.String("user_id", userID.String()))
			return nil, fmt.Errorf("failed to hash recovery code: %v", err)
		}
		hashedRecoveryCodes[i] = hashedCode
	}

	rowsAffected, err := s.Store.Enable2Fa(ctx, db.Enable2FaParams{
		ID:              user.ID,
		TwoFactorSecret: user.TwoFactorSecretTemp,
		RecoveryCodes:   hashedRecoveryCodes,
	})
	if err != nil || rowsAffected != 1 {
		s.Logger.LogError(ctx, "EnableTwoFA", "Database error enabling 2FA", err,
			zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to enable 2FA: %v", err)
	}

	s.Logger.LogInfo(ctx, "EnableTwoFA", "2FA enabled successfully",
		zap.String("user_id", userID.String()))
	return &Enable2FAResponse{
		RecoveryCodes: recoveryCodes,
	}, nil
}
