package api

import (
	"fmt"
	"net/http"

	"maicare_go/service/auth"

	"github.com/gin-gonic/gin"
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

	loginResult, err := server.businessService.AuthService.Login(auth.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		ClientIP:  ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res := SuccessResponse(LoginUserResponse{
		AccessToken:   loginResult.AccessToken,
		RefreshToken:  loginResult.RefreshToken,
		RequiresTwoFA: loginResult.RequiresTwoFA,
		TempToken:     loginResult.TempToken,
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
// @Description Refresh access token using refresh token“
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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}

	result, err := server.businessService.AuthService.RefreshToken(auth.RefreshTokenRequest{
		RefreshToken: req.Token,
	}, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("failed to refresh token: %v", err)))
		return
	}

	res := SuccessResponse(RefreshTokenResponse{
		AccessToken: result.AccessToken,
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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}

	loginResult, err := server.businessService.AuthService.VerifyTwoFAToken(auth.VerifyTwoFATokenRequest{
		ValidationCode: req.ValidationCode,
		TempToken:      req.TempToken,
	}, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("2FA verification failed: %v", err)))
		return
	}
	res := SuccessResponse(LoginUserResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
	}, "login successful")
	ctx.JSON(http.StatusOK, res)

}

// @Summary Logout user
// @Description Logout user and invalidate refresh token
// @Tags authentication
// @Produce json
// @Success 200 {object} Response[any] "Logout successful"
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /auth/logout [post]
func (server *Server) LogOutApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	err = server.businessService.AuthService.Logout(auth.LogoutRequest{
		PayloadID: payload.ID,
	}, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse[any](nil, "logout successful")
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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}
	err = server.businessService.AuthService.ChangePassword(auth.ChangePasswordRequest{
		UserID:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to change password")))
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
	setup2FAResult, err := server.businessService.AuthService.SetupTwoFA(auth.SetupTwoFARequest{
		UserID: userID,
	}, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to setup 2FA: %v", err)))
		return
	}

	res := SuccessResponse(Setup2FAResponse{
		QrCode: "data:image/png;base64," + setup2FAResult.QRCodeBase64,
		Secret: setup2FAResult.Secret,
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	userID := payload.UserId

	var req Enable2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}
	result, err := server.businessService.AuthService.EnableTwoFA(auth.EnableTwoFARequest{
		UserID:         userID,
		ValidationCode: req.ValidationCode,
	}, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to enable 2FA")))
		return
	}

	res := SuccessResponse(Enable2FAResponse{
		RecoveryCodes: result.RecoveryCodes,
	}, "2FA enabled successfully")
	ctx.JSON(http.StatusOK, res)
}
