package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	RoleType    string `json:"roleType"`
	IsActive    *bool  `json:"isActive"`
}

type RoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RoleType    string `json:"roleType"`
	IsActive    *bool  `json:"isActive"`
}

type RolePermissionUpdateRequest struct {
	PermissionIds []int `json:"permissionIds" binding:"required"`
}

func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, roles)
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_role_id", "Invalid role ID")
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), service.RoleCreateInput{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		RoleType:    req.RoleType,
		IsActive:    isActive,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Created(c, role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_role_id", "Invalid role ID")
		return
	}

	var req RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	role, err := h.roleService.UpdateRole(c.Request.Context(), id, service.RoleUpdateInput{
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		IsActive:    isActive,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, role)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_role_id", "Invalid role ID")
		return
	}

	err = h.roleService.DeleteRole(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Role deleted successfully"})
}

func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_role_id", "Invalid role ID")
		return
	}

	permissions, err := h.roleService.GetRolePermissions(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"roleId":      id,
		"permissions": permissions,
	})
}

func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_role_id", "Invalid role ID")
		return
	}

	var req RolePermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	err = h.roleService.UpdateRolePermissions(c.Request.Context(), id, req.PermissionIds)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"roleId":          id,
		"permissionCount": len(req.PermissionIds),
		"success":         true,
		"message":         "Permissions updated successfully",
	})
}

func (h *RoleHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, middlewares ...gin.HandlerFunc) {
	roles := router.Group("/roles")
	roles.Use(authMiddleware)
	if len(middlewares) > 0 {
		roles.Use(middlewares...)
	}
	{
		roles.GET("", h.GetRoles)
		roles.GET("/:id", h.GetRole)
		roles.POST("", h.CreateRole)
		roles.PUT("/:id", h.UpdateRole)
		roles.DELETE("/:id", h.DeleteRole)
		roles.GET("/:id/permissions", h.GetRolePermissions)
		roles.PUT("/:id/permissions", h.UpdateRolePermissions)
	}
}
