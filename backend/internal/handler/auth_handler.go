package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	TokenType        string `json:"tokenType"`
	ExpiresIn        int64  `json:"expiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordConfirm struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
	ConfirmPassword  string `json:"confirmPassword" binding:"required"`
}

type PasswordChangeRequest struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	if req.Username == "" && req.Email == "" {
		response.BadRequest(c, "missing_identifier", "Username or email is required")
		return
	}

	// Look up user
	var user *service.UserDTO
	var err error
	if req.Username != "" {
		user, err = h.authService.GetUserByUsername(c.Request.Context(), req.Username)
	} else {
		user, err = h.authService.GetUserByEmail(c.Request.Context(), req.Email)
	}

	if err != nil {
		response.Unauthorized(c, "invalid_credentials", "Invalid credentials")
		return
	}

	// Verify password
	if !h.authService.VerifyPassword(c.Request.Context(), req.Password, user.HashedPassword, user.Salt) {
		response.Unauthorized(c, "invalid_credentials", "Invalid credentials")
		return
	}

	// Check if user is active
	if !user.IsActive {
		response.Unauthorized(c, "account_disabled", "Account is disabled")
		return
	}

	// Generate tokens
	accessToken, err := h.authService.CreateAccessToken(c.Request.Context(), user.ID, user.Username, user.TokenVersion)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	refreshToken, err := h.authService.CreateRefreshToken(c.Request.Context(), user.ID, user.Username, user.TokenVersion)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        1800,
		RefreshExpiresIn: 604800,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	claims, err := h.authService.VerifyToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid_refresh_token", "Invalid refresh token")
		return
	}

	if claims.TokenType != "refresh" {
		response.Unauthorized(c, "invalid_token_type", "Invalid token type")
		return
	}

	// Verify TokenVersion against database before issuing new tokens
	user, err := h.authService.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Unauthorized(c, "user_not_found", "User not found")
		return
	}

	if claims.TokenVersion != user.TokenVersion {
		response.Unauthorized(c, "token_revoked", "Token has been revoked (password changed)")
		return
	}

	accessToken, err := h.authService.CreateAccessToken(c.Request.Context(), claims.UserID, claims.Username, claims.TokenVersion)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	refreshToken, err := h.authService.CreateRefreshToken(c.Request.Context(), claims.UserID, claims.Username, claims.TokenVersion)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        1800,
		RefreshExpiresIn: 604800,
	})
}

func (h *AuthHandler) RequestResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	code, expiresAt, err := h.authService.RequestResetPassword(c.Request.Context(), req.Email)
	if err != nil {
		// Don't reveal whether the email exists
		response.Success(c, gin.H{
			"message":   "If the email exists, a verification code has been sent",
			"expiresIn": 1800,
		})
		return
	}

	_ = code
	_ = expiresAt

	// TODO: Send code via email in production
	response.Success(c, gin.H{
		"message":   "Verification code sent",
		"expiresIn": 1800,
	})
}

func (h *AuthHandler) ConfirmResetPassword(c *gin.Context) {
	var req ResetPasswordConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		response.BadRequest(c, "password_mismatch", "Passwords do not match")
		return
	}

	if !h.authService.ValidatePasswordStrength(c.Request.Context(), req.NewPassword) {
		response.BadRequest(c, "weak_password", "Password must be at least 8 characters with letters and numbers")
		return
	}

	if err := h.authService.ConfirmResetPassword(c.Request.Context(), req.Email, req.VerificationCode, req.NewPassword); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid verification code") {
			response.BadRequest(c, "invalid_code", "Invalid verification code")
		} else if strings.Contains(errMsg, "expired") {
			response.BadRequest(c, "code_expired", "Verification code has expired")
		} else if strings.Contains(errMsg, "not found") {
			response.BadRequest(c, "user_not_found", "User not found")
		} else {
			response.InternalError(c, "reset_failed", "Password reset failed")
		}
		return
	}

	response.Success(c, gin.H{
		"message": "Password reset successfully",
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		response.BadRequest(c, "password_mismatch", "Passwords do not match")
		return
	}

	if !h.authService.ValidatePasswordStrength(c.Request.Context(), req.NewPassword) {
		response.BadRequest(c, "weak_password", "Password must be at least 8 characters with letters and numbers")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "not_authenticated", "User not authenticated")
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "old password is incorrect") {
			response.BadRequest(c, "incorrect_password", "Old password is incorrect")
		} else if strings.Contains(errMsg, "not found") {
			response.NotFound(c, "user_not_found", "User not found")
		} else if strings.Contains(errMsg, "strength") {
			response.BadRequest(c, "weak_password", "Password does not meet strength requirements")
		} else {
			response.InternalError(c, "password_change_failed", "Failed to change password")
		}
		return
	}

	response.Success(c, gin.H{
		"message": "Password changed successfully",
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "not_authenticated", "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "user_not_found", "User not found")
		return
	}

	response.Success(c, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"nickname":     user.Nickname,
		"userType":     user.UserType,
		"isActive":     user.IsActive,
		"tokenVersion": user.TokenVersion,
		"createdAt":    user.CreatedAt,
		"updatedAt":    user.UpdatedAt,
	})
}

func (h *AuthHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/reset/send-code", h.RequestResetPassword)
		auth.POST("/reset/confirm", h.ConfirmResetPassword)

		// Protected routes
		protected := auth.Group("")
		protected.Use(authMiddleware)
		{
			protected.POST("/change-password", h.ChangePassword)
			protected.GET("/me", h.GetCurrentUser)
		}
	}
}
