package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// LoginUserRequest represents the login request payload
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
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
// @Success 200 {object} LoginUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /token [post]
func (server *Server) Login(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error fetching user",
		})
		log.Printf("error getting user by username: %v", err)
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "incorrect password",
		})
		log.Printf("failed password check for user: %v", err)
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error generating access token",
		})
		log.Printf("failed to create access token for user : %v", err)
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error generating refresh token",
		})
		log.Printf("failed to create refresh token for user : %v", err)
		return
	}

	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: refreshPayload.ID, Valid: true},
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
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

	res := LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, res)
}
