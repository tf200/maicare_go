package auth

import (
	"context"
	"database/sql"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/token"
	"maicare_go/util"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// Methods and types for AuthService

type LoginRequest struct {
	Email     string
	Password  string
	ClientIP  string
	UserAgent string
}

type LoginResult struct {
	AccessToken   string
	RefreshToken  string
	RequiresTwoFA bool
	TempToken     string
}

func (s *authService) Login(req LoginRequest, ctx context.Context) (*LoginResult, error) {
	email := strings.ToLower(req.Email)

	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "Login", "Failed login attempt: user not found",
				zap.String("email", email), zap.String("client_ip", req.ClientIP), zap.String("user_agent", req.UserAgent))
			return nil, ErrInvalidCredentials
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Database error during login", zap.String("email", email),
			zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user")
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "Login", "Failed login attempt: incorrect password",
			zap.String("email", email), zap.String("client_ip", req.ClientIP),
			zap.String("user_agent", req.UserAgent))
		return nil, ErrInvalidCredentials
	}

	if user.TwoFactorEnabled {
		tempToken, _, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID,
			s.Config.TwoFATokenDuration, token.TwoFAToken)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create 2FA token",
				zap.String("email", email), zap.String("error", err.Error()))
			return nil, fmt.Errorf("failed to create 2FA token: %v", err)
		}
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "Login", "2FA required for user",
			zap.String("email", email), zap.String("client_ip", req.ClientIP),
			zap.String("user_agent", req.UserAgent))
		return &LoginResult{
			RequiresTwoFA: true,
			TempToken:     tempToken,
		}, nil
	}

	accessToken, _, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID,
		s.Config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create access token",
			zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create access token")
	}

	refreshToken, _, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID, s.Config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create refresh token",
			zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "Login", "User logged in successfully",
		zap.String("email", email), zap.String("client_ip", req.ClientIP),
		zap.String("user_agent", req.UserAgent))

	return &LoginResult{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		RequiresTwoFA: false,
		TempToken:     "",
	}, nil
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

func (s *authService) RefreshToken(req RefreshTokenRequest, ctx context.Context) (*RefreshTokenResult, error) {
	payload, err := s.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "RefreshToken", "Invalid refresh token",
			zap.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}

	session, err := s.Store.GetSessionByID(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "RefreshToken", "Session not found",
				zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))
			return nil, ErrSessionNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "RefreshToken", "Database error during session retrieval",
			zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()),
			zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	if session.IsBlocked {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "RefreshToken", "Blocked session attempt",
			zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}

	if session.UserID != payload.UserId {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "RefreshToken", "Session user mismatch",
			zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}

	if session.RefreshToken != req.RefreshToken {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "RefreshToken", "Refresh token mismatch",
			zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}
	if time.Now().After(session.ExpiresAt.Time) {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "RefreshToken", "Expired session attempt",
			zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}

	accessToken, _, err := s.TokenMaker.CreateToken(payload.UserId, payload.EmployeeID,
		s.Config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "RefreshToken", "Failed to create access token",
			zap.Int64("user_id", payload.UserId), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create access token")
	}

	result := &RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "RefreshToken", "Access token refreshed successfully",
		zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))

	return result, nil
}

type VerifyTwoFATokenRequest struct {
	ValidationCode string
	TempToken      string
}

func (s *authService) VerifyTwoFAToken(req VerifyTwoFATokenRequest, ctx context.Context) (*LoginResult, error) {
	tempPayload, err := s.TokenMaker.VerifyToken(req.TempToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "VerifyTwoFAToken", "Invalid temporary 2FA token",
			zap.String("error", err.Error()))
		return nil, ErrUnauthorized
	}

	user, err := s.Store.GetUserByID(context.Background(), tempPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "VerifyTwoFAToken", "User not found for 2FA",
				zap.Int64("user_id", tempPayload.UserId))
			return nil, ErrUserNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "VerifyTwoFAToken", "Database error during user retrieval",
			zap.Int64("user_id", tempPayload.UserId), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user")
	}

	if !user.TwoFactorEnabled || user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "VerifyTwoFAToken", "2FA not enabled for user",
			zap.Int64("user_id", user.ID))
		return nil, ErrUnauthorized
	}

	valid := totp.Validate(req.ValidationCode, *user.TwoFactorSecret)
	if !valid {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "VerifyTwoFAToken", "Invalid 2FA code",
			zap.Int64("user_id", user.ID))
		return nil, ErrUnauthorized
	}
	accessToken, _, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID, s.Config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "VerifyTwoFAToken", "Failed to create access token",
			zap.Int64("user_id", user.ID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create access token: %v", err)
	}

	refreshToken, _, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID, s.Config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "VerifyTwoFAToken", "Failed to create refresh token",
			zap.Int64("user_id", user.ID), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "VerifyTwoFAToken", "2FA verification successful, user logged in",
		zap.Int64("user_id", user.ID))

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type LogoutRequest struct {
	PayloadID uuid.UUID
}

func (s *authService) Logout(req LogoutRequest, ctx context.Context) error {
	err := s.Store.DeleteSession(ctx, req.PayloadID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "Logout", "Session not found during logout",
				zap.String("session_id", req.PayloadID.String()))
			return ErrSessionNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Logout", "Database error during session deletion",
			zap.String("session_id", req.PayloadID.String()), zap.String("error", err.Error()))
		return fmt.Errorf("failed to delete session: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "Logout", "User logged out successfully",
		zap.String("session_id", req.PayloadID.String()))

	return nil
}

type ChangePasswordRequest struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

func (s *authService) ChangePassword(req ChangePasswordRequest, ctx context.Context) error {
	user, err := s.Store.GetUserByID(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "ChangePassword", "User not found during password change",
				zap.Int64("user_id", req.UserID))
			return ErrUserNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ChangePassword", "Database error during user retrieval",
			zap.Int64("user_id", req.UserID), zap.String("error", err.Error()))
		return fmt.Errorf("failed to get user: %v", err)
	}

	err = util.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "ChangePassword", "Incorrect old password during password change",
			zap.Int64("user_id", req.UserID))
		return ErrInvalidCredentials
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ChangePassword", "Error hashing new password",
			zap.Int64("user_id", req.UserID), zap.String("error", err.Error()))
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	err = s.Store.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:       req.UserID,
		Password: hashedPassword,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ChangePassword", "Database error updating password",
			zap.Int64("user_id", req.UserID), zap.String("error", err.Error()))
		return fmt.Errorf("failed to update password: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ChangePassword", "Password changed successfully",
		zap.Int64("user_id", req.UserID))

	return nil
}
