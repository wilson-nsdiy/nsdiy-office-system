package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type PermissionHandler struct {
	permissionService *service.PermissionService
}

func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

type PermissionCreateRequest struct {
	Pid          *int   `json:"pid"`
	Name         string `json:"name" binding:"required"`
	ResourceType string `json:"resourceType" binding:"required"`
	ResourcePath string `json:"resourcePath" binding:"required"`
	HTTPMethod   string `json:"httpMethod"`
	Description  string `json:"description"`
	IsActive     *bool  `json:"isActive"`
}

type PermissionUpdateRequest struct {
	Pid          *int   `json:"pid"`
	Name         string `json:"name"`
	ResourceType string `json:"resourceType"`
	ResourcePath string `json:"resourcePath"`
	HTTPMethod   string `json:"httpMethod"`
	Description  string `json:"description"`
	IsActive     *bool  `json:"isActive"`
}

func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	resourceType := c.Query("resourceType")
	keyword := c.Query("keyword")

	permissions, err := h.permissionService.GetPermissions(c.Request.Context(), resourceType, keyword)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, permissions)
}

func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.permissionService.GetAllPermissions(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, permissions)
}

func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_permission_id", "Invalid permission ID")
		return
	}

	perm, err := h.permissionService.GetPermission(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, perm)
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req PermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	perm, err := h.permissionService.CreatePermission(c.Request.Context(), service.PermissionCreateInput{
		Pid:          req.Pid,
		Name:         req.Name,
		ResourceType: req.ResourceType,
		ResourcePath: req.ResourcePath,
		HTTPMethod:   req.HTTPMethod,
		Description:  req.Description,
		IsActive:     isActive,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, perm)
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_permission_id", "Invalid permission ID")
		return
	}

	var req PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	perm, err := h.permissionService.UpdatePermission(c.Request.Context(), id, service.PermissionUpdateInput{
		Pid:          req.Pid,
		Name:         req.Name,
		ResourceType: req.ResourceType,
		ResourcePath: req.ResourcePath,
		HTTPMethod:   req.HTTPMethod,
		Description:  req.Description,
		IsActive:     isActive,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, perm)
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_permission_id", "Invalid permission ID")
		return
	}

	err = h.permissionService.DeletePermission(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Permission deleted successfully"})
}

func (h *PermissionHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, middlewares ...gin.HandlerFunc) {
	permissions := router.Group("/permissions")
	permissions.Use(authMiddleware)
	if len(middlewares) > 0 {
		permissions.Use(middlewares...)
	}
	{
		permissions.GET("", h.GetPermissions)
		permissions.GET("/all", h.GetAllPermissions)
		permissions.GET("/:id", h.GetPermission)
		permissions.POST("", h.CreatePermission)
		permissions.PUT("/:id", h.UpdatePermission)
		permissions.DELETE("/:id", h.DeletePermission)
	}
}
