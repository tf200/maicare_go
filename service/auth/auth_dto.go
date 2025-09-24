package auth

// LoginUserRequest represents the login request payload
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"testemail@gmail.com"`
	Password string `json:"password" binding:"required" example:"t2aha000"`
}

// LoginUserRequest represents the login request payload
type LoginUserResponse struct {
	RefreshToken  string `json:"refresh" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	AccessToken   string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RequiresTwoFA bool   `json:"requires_2fa" example:"false"`
	TempToken     string `json:"temp_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Verify2FARequest represents the verify 2FA request payload
type Verify2FARequest struct {
	ValidationCode string `json:"validation_code" binding:"required"`
	TempToken      string `json:"temp_token" binding:"required"`
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Setup2FARequest represents the setup 2FA request payload
type Setup2FAResponse struct {
	QrCode string `json:"qr_code_base64" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Secret string `json:"secret" example:"JBSWY3DPEHPK3PXP"`
}
