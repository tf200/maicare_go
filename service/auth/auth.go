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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// Methods and types for AuthService

func (s *authService) Login(req LoginUserRequest, clientIP string,
	userAgent string, ctx context.Context) (*LoginUserResponse, error) {

	email := strings.ToLower(req.Email)

	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "Login", "Failed login attempt: user not found",
				zap.String("email", email), zap.String("client_ip", clientIP), zap.String("user_agent", userAgent))
			return nil, ErrInvalidCredentials
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Database error during login", zap.String("email", email),
			zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user")
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "Login", "Failed login attempt: incorrect password",
			zap.String("email", email), zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent))
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
			zap.String("email", email), zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent))
		return &LoginUserResponse{
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

	refreshToken, payload, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID, s.Config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create refresh token",
			zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	session, err := s.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payload.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: payload.ExpiresAt, Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UserID:       payload.UserId,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "Login", "Database error during session creation",
			zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "Login", "User logged in successfully",
		zap.String("email", email), zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent), zap.String("session_id", session.ID.String()))

	return &LoginUserResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		RequiresTwoFA: false,
		TempToken:     "",
	}, nil
}

func (s *authService) RefreshToken(req RefreshTokenRequest, ctx context.Context) (*RefreshTokenResponse, error) {
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

	result := &RefreshTokenResponse{
		AccessToken: accessToken,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "RefreshToken", "Access token refreshed successfully",
		zap.Int64("user_id", payload.UserId), zap.String("session_id", payload.ID.String()))

	return result, nil
}

func (s *authService) VerifyTwoFAToken(req Verify2FARequest, ctx context.Context) (*LoginUserResponse, error) {
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

	return &LoginUserResponse{
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

func (s *authService) ChangePassword(req ChangePasswordRequest, userID int64, ctx context.Context) error {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "ChangePassword", "User not found during password change",
				zap.Int64("user_id", userID))
			return ErrUserNotFound
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ChangePassword", "Database error during user retrieval",
			zap.Int64("user_id", userID), zap.String("error", err.Error()))
		return fmt.Errorf("failed to get user: %v", err)
	}

	err = util.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelWarn, "ChangePassword", "Incorrect old password during password change",
			zap.Int64("user_id", userID))
		return ErrInvalidCredentials
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ChangePassword", "Error hashing new password",
			zap.Int64("user_id", userID), zap.String("error", err.Error()))
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	err = s.Store.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:       userID,
		Password: hashedPassword,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ChangePassword", "Database error updating password",
			zap.Int64("user_id", userID), zap.String("error", err.Error()))
		return fmt.Errorf("failed to update password: %v", err)
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ChangePassword", "Password changed successfully",
		zap.Int64("user_id", userID))

	return nil
}
