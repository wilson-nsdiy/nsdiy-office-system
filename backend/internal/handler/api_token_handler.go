package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type ApiTokenHandler struct {
	apiTokenService *service.ApiTokenService
}

func NewApiTokenHandler(apiTokenService *service.ApiTokenService) *ApiTokenHandler {
	return &ApiTokenHandler{apiTokenService: apiTokenService}
}

type ApiTokenCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt string `json:"expiresAt"`
}

type ApiTokenUpdateRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (h *ApiTokenHandler) ListTokens(c *gin.Context) {
	userID := middleware.GetUserID(c)
	keyword := c.Query("keyword")
	status := c.Query("status")
	page, pageSize := response.ParsePagination(c)

	result, err := h.apiTokenService.ListTokens(c.Request.Context(), userID, keyword, status, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *ApiTokenHandler) GetToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_token_id", "Invalid token ID")
		return
	}

	userID := middleware.GetUserID(c)

	token, err := h.apiTokenService.GetToken(c.Request.Context(), id, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, token)
}

func (h *ApiTokenHandler) CreateToken(c *gin.Context) {
	var req ApiTokenCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	result, err := h.apiTokenService.CreateToken(c.Request.Context(), userID, service.ApiTokenCreateInput{
		Name:      req.Name,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Created(c, gin.H{
		"id":          result.Token.ID,
		"name":        result.Token.Name,
		"token":       result.RawToken,
		"tokenPrefix": result.Token.TokenPrefix,
		"status":      result.Token.Status,
		"message":     "Token created successfully. Please save the token as it will not be shown again.",
	})
}

func (h *ApiTokenHandler) UpdateToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_token_id", "Invalid token ID")
		return
	}

	userID := middleware.GetUserID(c)

	var req ApiTokenUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	token, err := h.apiTokenService.UpdateToken(c.Request.Context(), id, userID, service.ApiTokenUpdateInput{
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, token)
}

func (h *ApiTokenHandler) DeleteToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_token_id", "Invalid token ID")
		return
	}

	userID := middleware.GetUserID(c)

	err = h.apiTokenService.DeleteToken(c.Request.Context(), id, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Token deleted successfully"})
}

func (h *ApiTokenHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	apiTokens := router.Group("/api-tokens")
	apiTokens.Use(authMiddleware)
	{
		apiTokens.GET("", h.ListTokens)
		apiTokens.GET("/:id", h.GetToken)
		apiTokens.POST("", h.CreateToken)
		apiTokens.PUT("/:id", h.UpdateToken)
		apiTokens.DELETE("/:id", h.DeleteToken)
	}
}
