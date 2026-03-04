package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

func (s *authService) Login(req LoginUserRequest, clientIP string,
	userAgent string, ctx context.Context,
) (*LoginUserResponse, error) {
	email := strings.ToLower(req.Email)

	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogError(ctx, "Login", "Failed login attempt: user not found", nil,
				zap.String("email", email), zap.String("client_ip", clientIP), zap.String("user_agent", userAgent))
			return nil, ErrInvalidCredentials
		}
		s.Logger.LogError(ctx, "Login", "Database error during login", err, zap.String("email", email))
		return nil, fmt.Errorf("failed to get user")
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		s.Logger.LogError(ctx, "Login", "Failed login attempt: incorrect password", nil,
			zap.String("email", email), zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent))
		return nil, ErrInvalidCredentials
	}

	if user.TwoFactorEnabled {
		tempToken, _, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID,
			s.Config.TwoFATokenDuration, token.TwoFAToken)
		if err != nil {
			s.Logger.LogError(ctx, "Login", "Failed to create 2FA token", err, zap.String("email", email))
			return nil, fmt.Errorf("failed to create 2FA token: %v", err)
		}
		s.Logger.LogInfo(ctx, "Login", "2FA required for user",
			zap.String("email", email), zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent))
		return &LoginUserResponse{
			RequiresTwoFA: true,
			TempToken:     tempToken,
		}, nil
	}

	refreshToken, payload, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID, s.Config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		s.Logger.LogError(ctx, "Login", "Failed to create refresh token", err, zap.String("email", email))
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	accessToken, _, err := s.TokenMaker.CreateTokenWithSessionID(user.ID, user.EmployeeID,
		s.Config.AccessTokenDuration, token.AccessToken, payload.ID)
	if err != nil {
		s.Logger.LogError(ctx, "Login", "Failed to create access token", err, zap.String("email", email))
		return nil, fmt.Errorf("failed to create access token")
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
		s.Logger.LogError(ctx, "Login", "Database error during session creation", err, zap.String("email", email))
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	s.Logger.LogInfo(ctx, "Login", "User logged in successfully",
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
		s.Logger.LogWarn(ctx, "RefreshToken", "Invalid refresh token", zap.Error(err))
		return nil, ErrInvalidCredentials
	}

	session, err := s.Store.GetSessionByID(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogWarn(ctx, "RefreshToken", "Session not found",
				zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))
			return nil, ErrSessionNotFound
		}
		s.Logger.LogError(ctx, "RefreshToken", "Database error during session retrieval", err,
			zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	if session.IsBlocked {
		s.Logger.LogWarn(ctx, "RefreshToken", "Blocked session attempt",
			zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}

	if session.UserID != payload.UserId {
		s.Logger.LogWarn(ctx, "RefreshToken", "Session user mismatch",
			zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}

	if session.RefreshToken != req.RefreshToken {
		s.Logger.LogWarn(ctx, "RefreshToken", "Refresh token mismatch",
			zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}
	if time.Now().After(session.ExpiresAt.Time) {
		s.Logger.LogWarn(ctx, "RefreshToken", "Expired session attempt",
			zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))
		return nil, ErrUnauthorized
	}

	accessToken, _, err := s.TokenMaker.CreateTokenWithSessionID(payload.UserId, payload.EmployeeID,
		s.Config.AccessTokenDuration, token.AccessToken, payload.ID)
	if err != nil {
		s.Logger.LogError(ctx, "RefreshToken", "Failed to create access token", err,
			zap.String("user_id", payload.UserId.String()))
		return nil, fmt.Errorf("failed to create access token")
	}

	result := &RefreshTokenResponse{
		AccessToken: accessToken,
	}

	s.Logger.LogInfo(ctx, "RefreshToken", "Access token refreshed successfully",
		zap.String("user_id", payload.UserId.String()), zap.String("session_id", payload.ID.String()))

	return result, nil
}

func (s *authService) VerifyTwoFAToken(req Verify2FARequest, clientIP string, userAgent string, ctx context.Context) (*LoginUserResponse, error) {
	tempPayload, err := s.TokenMaker.VerifyToken(req.TempToken)
	if err != nil {
		s.Logger.LogWarn(ctx, "VerifyTwoFAToken", "Invalid temporary 2FA token", zap.Error(err))
		return nil, ErrUnauthorized
	}

	user, err := s.Store.GetUserByID(ctx, tempPayload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogWarn(ctx, "VerifyTwoFAToken", "User not found for 2FA",
				zap.String("user_id", tempPayload.UserId.String()))
			return nil, ErrUserNotFound
		}
		s.Logger.LogError(ctx, "VerifyTwoFAToken", "Database error during user retrieval", err,
			zap.String("user_id", tempPayload.UserId.String()))
		return nil, fmt.Errorf("failed to get user")
	}

	if !user.TwoFactorEnabled || user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		s.Logger.LogWarn(ctx, "VerifyTwoFAToken", "2FA not enabled for user",
			zap.String("user_id", user.ID.String()))
		return nil, ErrUnauthorized
	}

	valid := totp.Validate(req.ValidationCode, *user.TwoFactorSecret)
	if !valid {
		s.Logger.LogWarn(ctx, "VerifyTwoFAToken", "Invalid 2FA code",
			zap.String("user_id", user.ID.String()))
		return nil, ErrUnauthorized
	}
	refreshToken, payload, err := s.TokenMaker.CreateToken(user.ID, user.EmployeeID, s.Config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		s.Logger.LogError(ctx, "VerifyTwoFAToken", "Failed to create refresh token", err,
			zap.String("user_id", user.ID.String()))
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	accessToken, _, err := s.TokenMaker.CreateTokenWithSessionID(user.ID, user.EmployeeID, s.Config.AccessTokenDuration, token.AccessToken, payload.ID)
	if err != nil {
		s.Logger.LogError(ctx, "VerifyTwoFAToken", "Failed to create access token", err,
			zap.String("user_id", user.ID.String()))
		return nil, fmt.Errorf("failed to create access token: %v", err)
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
		s.Logger.LogError(ctx, "VerifyTwoFAToken", "Database error during session creation", err,
			zap.String("user_id", user.ID.String()))
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	s.Logger.LogInfo(ctx, "VerifyTwoFAToken", "2FA verification successful, user logged in",
		zap.String("user_id", user.ID.String()), zap.String("session_id", session.ID.String()))

	return &LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type LogoutRequest struct {
	SessionID uuid.UUID
}

func (s *authService) Logout(req LogoutRequest, ctx context.Context) error {
	err := s.Store.DeleteSession(ctx, req.SessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogWarn(ctx, "Logout", "Session not found during logout",
				zap.String("session_id", req.SessionID.String()))
			return ErrSessionNotFound
		}
		s.Logger.LogError(ctx, "Logout", "Database error during session deletion", err,
			zap.String("session_id", req.SessionID.String()))
		return fmt.Errorf("failed to delete session: %v", err)
	}

	s.Logger.LogInfo(ctx, "Logout", "User logged out successfully",
		zap.String("session_id", req.SessionID.String()))

	return nil
}

func (s *authService) ChangePassword(req ChangePasswordRequest, userID uuid.UUID, ctx context.Context) error {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logger.LogWarn(ctx, "ChangePassword", "User not found during password change",
				zap.String("user_id", userID.String()))
			return ErrUserNotFound
		}
		s.Logger.LogError(ctx, "ChangePassword", "Database error during user retrieval", err,
			zap.String("user_id", userID.String()))
		return fmt.Errorf("failed to get user: %v", err)
	}

	err = util.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		s.Logger.LogWarn(ctx, "ChangePassword", "Incorrect old password during password change",
			zap.String("user_id", userID.String()))
		return ErrInvalidCredentials
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		s.Logger.LogError(ctx, "ChangePassword", "Error hashing new password", err,
			zap.String("user_id", userID.String()))
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	err = s.Store.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:       userID,
		Password: hashedPassword,
	})
	if err != nil {
		s.Logger.LogError(ctx, "ChangePassword", "Database error updating password", err,
			zap.String("user_id", userID.String()))
		return fmt.Errorf("failed to update password: %v", err)
	}

	s.Logger.LogInfo(ctx, "ChangePassword", "Password changed successfully",
		zap.String("user_id", userID.String()))

	return nil
}
