package api

import (
	"fmt"
	"net/http"

	"maicare_go/service/auth"

	"github.com/gin-gonic/gin"
)

// @Summary Generate authentication tokens
// @Description Authenticate user and return access and refresh tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body auth.LoginUserRequest true "Login credentials"
// @Success 200 {object} Response[auth.LoginUserResponse] "Successfully authenticated"
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - Account status issue"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /auth/token [post]
// @Security -
func (server *Server) Login(ctx *gin.Context) {
	var req auth.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loginResult, err := server.businessService.AuthService.Login(req, ctx.ClientIP(), ctx.Request.UserAgent(), ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res := SuccessResponse(loginResult, "login successful")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Refresh access token

// @Description Refresh access token using refresh token“
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body auth.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} Response[auth.RefreshTokenResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /auth/refresh [post]
// @Security -
func (server *Server) RefreshToken(ctx *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}
	result, err := server.businessService.AuthService.RefreshToken(req, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("failed to refresh token: %v", err)))
		return
	}

	res := SuccessResponse(result, "access token refreshed successfully")

	ctx.JSON(http.StatusOK, res)
}

// Verify2FAHandler verifies the 2FA code and generates access and refresh tokens
// @Summary Verify 2FA code
// @Description Verify 2FA code and generate access and refresh tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body auth.Verify2FARequest true "Verify 2FA request"
// @Success 200 {object} Response[auth.LoginUserResponse] "2FA verification successful"
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /auth/verify_2fa [post]
// @Security -
func (server *Server) Verify2FAHandler(ctx *gin.Context) {
	var req auth.Verify2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}

	loginResult, err := server.businessService.AuthService.VerifyTwoFAToken(req, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("2FA verification failed")))
		return
	}
	res := SuccessResponse(loginResult, "login successful")
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

// @Summary Change user password
// @Description Change user password
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body auth.ChangePasswordRequest true "Change password request"
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	userID := payload.UserId

	var req auth.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}
	err = server.businessService.AuthService.ChangePassword(req, userID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to change password")))
		return
	}
	res := SuccessResponse[any](nil, "password changed successfully")
	ctx.JSON(http.StatusOK, res)

}

// @Summary Setup 2FA
// @Description Setup 2FA for user
// @Tags authentication
// @Produce json
// @Success 200 {object} Response[auth.Setup2FAResponse] "2FA setup successful"
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
	setup2FAResult, err := server.businessService.AuthService.SetupTwoFA(userID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to setup 2FA: %v", err)))
		return
	}

	res := SuccessResponse(setup2FAResult, "2FA setup successful")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Enable 2FA
// @Description Enable 2FA for user
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body auth.Enable2FARequest true "Enable 2FA request"
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

	var req auth.Enable2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}
	result, err := server.businessService.AuthService.EnableTwoFA(req, userID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to enable 2FA")))
		return
	}

	res := SuccessResponse(result, "2FA enabled successfully")
	ctx.JSON(http.StatusOK, res)
}
