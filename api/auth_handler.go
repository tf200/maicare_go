package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// LoginUserRequest represents the login request payload
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required" example:"testemail@gmail.com"`
	Password string `json:"password" binding:"required" example:"t2aha000"`
}

// LoginUserResponse represents the login response
type LoginUserResponse struct {
	RefreshToken string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	AccessToken  string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
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

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		log.Printf("error getting user by username: %v", err)
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		log.Printf("failed password check for user: %v", err)
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
