package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// LoginUserRequest represents the login request payload
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required" example:"testemail@gmail.com"`
	Password string `json:"password" binding:"required" example:"t2aha000"`
}

// LoginUserResponse represents the login response
type LoginUserResponse struct {
	RefreshToken  string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	AccessToken   string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RequiresTwoFA bool   `json:"requires_2fa" example:"false"`
	TempToken     string `json:"temp_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// @Summary Generate authentication tokens
// @Description Authenticate user and return access and refresh tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginUserRequest true "Login credentials"
// @Success 200 {object} Response[LoginUserResponse] "Successfully authenticated"
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - Account status issue"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /auth/token [post]
// @Security -
func (server *Server) Login(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert email to lowercase
	lowerCaseEmail := strings.ToLower(req.Email)

	user, err := server.store.GetUserByEmail(ctx, lowerCaseEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("invalid password or email for user")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		log.Printf("error getting user by username: %v", err)
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid password or email for user")))
		log.Printf("failed password check for user: %v", err)
		return
	}

	if user.TwoFactorEnabled {
		tempToken, _, err := server.tokenMaker.CreateToken(user.ID, int32(user.RoleID), server.config.TwoFATokenDuration, token.TwoFAToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		res := SuccessResponse(LoginUserResponse{
			RequiresTwoFA: true,
			TempToken:     tempToken,
		}, "2FA required")
		ctx.JSON(http.StatusOK, res)
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, int32(user.RoleID), server.config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, int32(user.RoleID), server.config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiresAt, Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: refreshPayload.IssuedAt, Valid: true},
		UserID:       user.ID,
	})
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			case "23503": // foreign_key_violation
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "login successful")

	ctx.JSON(http.StatusOK, res)
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} Response[RefreshTokenResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /auth/refresh [post]
// @Security -
func (server *Server) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	session, err := server.store.GetSessionByID(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("blocked session")))
		return
	}

	if session.RefreshToken != req.Token {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("mismatched session token")))
		return
	}
	if time.Now().After(session.ExpiresAt.Time) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("expired session")))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(payload.UserId, payload.RoleID, server.config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error generating access token",
		})
		return
	}

	res := SuccessResponse(RefreshTokenResponse{
		AccessToken: accessToken,
	}, "access token refreshed successfully")

	ctx.JSON(http.StatusOK, res)
}

// Verify2FARequest represents the verify 2FA request payload
type Verify2FARequest struct {
	ValidationCode string `json:"validation_code" binding:"required"`
	TempToken      string `json:"temp_token" binding:"required"`
}

// Verify2FAHandler verifies the 2FA code and generates access and refresh tokens
// @Summary Verify 2FA code
// @Description Verify 2FA code and generate access and refresh tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body Verify2FARequest true "Verify 2FA request"
// @Success 200 {object} Response[LoginUserResponse] "2FA verification successful"
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /auth/verify_2fa [post]
// @Security -
func (server *Server) Verify2FAHandler(ctx *gin.Context) {
	var req Verify2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tempPayload, err := server.tokenMaker.VerifyToken(req.TempToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByID(ctx, tempPayload.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !user.TwoFactorEnabled || user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("2FA not enabled or secret not set")))
		return
	}
	valid := totp.Validate(req.ValidationCode, *user.TwoFactorSecret)
	if !valid {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid validation code")))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, int32(user.RoleID), server.config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, int32(user.RoleID), server.config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshPayload.ExpiresAt, Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: refreshPayload.IssuedAt, Valid: true},
		UserID:       user.ID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "login successful")
	ctx.JSON(http.StatusOK, res)

}

// LogoutRequest represents the logout request payload
type LogoutResponse struct {
	Message string `json:"message" example:"logout successful"`
}

// @Summary Logout user
// @Description Logout user and invalidate refresh token
// @Tags authentication
// @Produce json
// @Success 200 {object} Response[LogoutResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /auth/logout [post]
func (server *Server) LogOutApi(ctx *gin.Context) {
	payload, exist := ctx.Get(authorizationPayloadKey)
	if !exist {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authorization payload not found")))
		return
	}
	sessionPayload, ok := payload.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("invalid authorization payload type")))
		return
	}

	err := server.store.DeleteSession(ctx, sessionPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(LogoutResponse{
		Message: "logout successful",
	}, "logout successful")
	ctx.JSON(http.StatusOK, res)
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// @Summary Change user password
// @Description Change user password
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} Response[any] "Password changed successfully"
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - Password change issue"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /auth/change_password [post]
func (server *Server) ChangePasswordApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userID := payload.UserId

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		log.Printf("failed password check for user: %v", err)
		return
	}
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:       user.ID,
		Password: hashedPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			case "23503": // foreign_key_violation
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse[any](nil, "password changed successfully")
	ctx.JSON(http.StatusOK, res)

}

// Setup2FARequest represents the setup 2FA request payload
type Setup2FAResponse struct {
	QrCode string `json:"qr_code_base64" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Secret string `json:"secret" example:"JBSWY3DPEHPK3PXP"`
}

// @Summary Setup 2FA
// @Description Setup 2FA for user
// @Tags authentication
// @Produce json
// @Success 200 {object} Response[Setup2FAResponse] "2FA setup successful"
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - 2FA setup issue"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /auth/setup_2fa [post]
func (server *Server) Setup2FAHandler(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))

		return
	}

	userID := payload.UserId

	user, err := server.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Maicare",
		AccountName: user.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	secretKey := key.Secret()

	err = server.store.CreateTemp2FaSecret(ctx, db.CreateTemp2FaSecretParams{
		ID:                  user.ID,
		TwoFactorSecretTemp: &secretKey,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	res := SuccessResponse(Setup2FAResponse{
		QrCode: "data:image/png;base64," + qrCodeBase64,
		Secret: secretKey,
	}, "2FA setup successful")
	ctx.JSON(http.StatusOK, res)
}

// Enable2FARequest represents the enable 2FA request payload
type Enable2FARequest struct {
	ValidationCode string `json:"validation_code" binding:"required"`
}

type Enable2FAResponse struct {
	RecoveryCodes []string `json:"recovery_codes" example:"[\"code1\", \"code2\"]"`
}

// @Summary Enable 2FA
// @Description Enable 2FA for user
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body Enable2FARequest true "Enable 2FA request"
// @Success 200 {object} Response[any] "2FA enabled successfully"
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - 2FA enable issue"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /auth/enable_2fa [post]
func (server *Server) Enable2FAHandler(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userID := payload.UserId

	var req Enable2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.TwoFactorSecretTemp == nil || *user.TwoFactorSecretTemp == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("2FA setup not initiated")))
		return
	}

	valid := totp.Validate(req.ValidationCode, *user.TwoFactorSecretTemp)
	if !valid {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid validation code")))
		return
	}

	recoveryCodes := util.GenerateRecoveryCodes(10)

	hashedRecoveryCodes := make([]string, len(recoveryCodes))
	for i, code := range recoveryCodes {
		hashedCode, err := util.HashPassword(code)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		hashedRecoveryCodes[i] = hashedCode
	}

	err = server.store.Enable2Fa(ctx, db.Enable2FaParams{
		ID:              user.ID,
		TwoFactorSecret: user.TwoFactorSecretTemp,
		RecoveryCodes:   hashedRecoveryCodes,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(Enable2FAResponse{
		RecoveryCodes: recoveryCodes,
	}, "2FA enabled successfully")
	ctx.JSON(http.StatusOK, res)
}
