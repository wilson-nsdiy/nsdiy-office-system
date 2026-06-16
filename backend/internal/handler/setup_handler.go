package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
	"oa-nsdiy/backend/internal/setup"
)

// installMutex prevents concurrent installation attempts
var installMutex sync.Mutex

type SetupHandler struct {
	authService *service.AuthService
}

func NewSetupHandler(authService *service.AuthService) *SetupHandler {
	return &SetupHandler{authService: authService}
}

// needsSetupDoubleCheck verifies installation state using both file lock and database.
// This prevents bypass if the .installed file is deleted.
func (h *SetupHandler) needsSetupDoubleCheck(ctx context.Context) (bool, error) {
	if !setup.NeedsSetup() {
		return false, nil
	}

	hasUser, err := h.authService.HasAnyUser(ctx)
	if err != nil {
		return false, err
	}

	if hasUser {
		if err := setup.CreateInstallLock(); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

type SetupStatusResponse struct {
	NeedsSetup bool `json:"needsSetup"`
}

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Nickname string `json:"nickname" binding:"required,min=1,max=100"`
}

// Status checks if the system needs initial setup
func (h *SetupHandler) Status(c *gin.Context) {
	needsSetup, err := h.needsSetupDoubleCheck(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, SetupStatusResponse{
		NeedsSetup: needsSetup,
	})
}

// CreateAdmin creates the initial admin user
func (h *SetupHandler) CreateAdmin(c *gin.Context) {
	installMutex.Lock()
	defer installMutex.Unlock()

	// Double-check: prevent re-installation using both file and database
	needsSetup, err := h.needsSetupDoubleCheck(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if !needsSetup {
		response.ErrorWithDetails(c, http.StatusForbidden, "System is already installed", "already_installed", nil)
		return
	}

	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	// Create admin user
	salt, err := h.authService.GenerateSalt(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	hashedPassword, err := h.authService.HashPassword(c.Request.Context(), req.Password, salt)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	input := &service.CreateUserInput{
		Username:       req.Username,
		Email:          req.Email,
		Nickname:       &req.Nickname,
		Salt:           salt,
		HashedPassword: hashedPassword,
		UserType:       "HUMAN",
		IsActive:       true,
	}

	if err := h.authService.CreateUser(c.Request.Context(), input); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Create install lock file to prevent re-installation
	if err := setup.CreateInstallLock(); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"message":  "Admin user created successfully",
		"username": req.Username,
	})
}

// SetupGuard returns a middleware that blocks setup endpoints if already installed.
// Uses double-check verification with both file lock and database.
func (h *SetupHandler) SetupGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		needsSetup, err := h.needsSetupDoubleCheck(c.Request.Context())
		if err != nil {
			response.ErrorFrom(c, err)
			c.Abort()
			return
		}
		if !needsSetup {
			response.ErrorWithDetails(c, http.StatusForbidden, "Setup is not allowed: system is already installed", "already_installed", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *SetupHandler) SetupRoutes(api *gin.RouterGroup) {
	setup := api.Group("/setup")
	{
		// Status endpoint is always accessible
		setup.GET("/status", h.Status)

		// Modification endpoints are protected by guard
		protected := setup.Group("")
		protected.Use(h.SetupGuard())
		{
			protected.POST("/admin", h.CreateAdmin)
		}
	}
}
