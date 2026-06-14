package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type MediaHandler struct {
	mediaAccountService *service.MediaAccountService
	mediaContentService *service.MediaContentService
}

func NewMediaHandler(mediaAccountService *service.MediaAccountService, mediaContentService *service.MediaContentService) *MediaHandler {
	return &MediaHandler{
		mediaAccountService: mediaAccountService,
		mediaContentService: mediaContentService,
	}
}

type MediaAccountCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	Platform  string `json:"platform" binding:"required"`
	AccountId string `json:"accountId" binding:"required"`
	Avatar    string `json:"avatar"`
}

type MediaAccountUpdateRequest struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	TokenExpiresAt string `json:"tokenExpiresAt"`
}

type MediaContentCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	CoverImage  string `json:"coverImage"`
	Platform    string `json:"platform" binding:"required"`
	AccountId   *int   `json:"accountId"`
	Status      string `json:"status"`
	PublishTime string `json:"publishTime"`
}

type MediaContentUpdateRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	CoverImage  string `json:"coverImage"`
	Status      string `json:"status"`
	Views       *int   `json:"views"`
	Likes       *int   `json:"likes"`
	Comments    *int   `json:"comments"`
	Shares      *int   `json:"shares"`
	PublishTime string `json:"publishTime"`
}

func (h *MediaHandler) ListAccounts(c *gin.Context) {
	keyword := c.Query("keyword")
	platform := c.Query("platform")
	status := c.Query("status")
	page, pageSize := response.ParsePagination(c)

	result, err := h.mediaAccountService.ListAccounts(c.Request.Context(), keyword, platform, status, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *MediaHandler) GetAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_account_id", "Invalid account ID")
		return
	}

	account, err := h.mediaAccountService.GetAccount(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, account)
}

func (h *MediaHandler) CreateAccount(c *gin.Context) {
	var req MediaAccountCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	account, err := h.mediaAccountService.CreateAccount(c.Request.Context(), service.MediaAccountCreateInput{
		Name:      req.Name,
		Platform:  req.Platform,
		AccountId: req.AccountId,
		Avatar:    req.Avatar,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, account)
}

func (h *MediaHandler) UpdateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_account_id", "Invalid account ID")
		return
	}

	var req MediaAccountUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	account, err := h.mediaAccountService.UpdateAccount(c.Request.Context(), id, service.MediaAccountUpdateInput{
		Name:         req.Name,
		Status:       req.Status,
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, account)
}

func (h *MediaHandler) DeleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_account_id", "Invalid account ID")
		return
	}

	err = h.mediaAccountService.DeleteAccount(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Account deleted successfully"})
}

func (h *MediaHandler) ListContents(c *gin.Context) {
	keyword := c.Query("keyword")
	platform := c.Query("platform")
	status := c.Query("status")

	var accountID *int
	if aid := c.Query("accountId"); aid != "" {
		id, _ := strconv.Atoi(aid)
		accountID = &id
	}

	page, pageSize := response.ParsePagination(c)

	result, err := h.mediaContentService.ListContents(c.Request.Context(), keyword, platform, status, accountID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *MediaHandler) GetContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_content_id", "Invalid content ID")
		return
	}

	content, err := h.mediaContentService.GetContent(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, content)
}

func (h *MediaHandler) CreateContent(c *gin.Context) {
	var req MediaContentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	content, err := h.mediaContentService.CreateContent(c.Request.Context(), service.MediaContentCreateInput{
		Title:       req.Title,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		Platform:    req.Platform,
		AccountId:   req.AccountId,
		Status:      req.Status,
		PublishTime: req.PublishTime,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, content)
}

func (h *MediaHandler) UpdateContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_content_id", "Invalid content ID")
		return
	}

	var req MediaContentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	content, err := h.mediaContentService.UpdateContent(c.Request.Context(), id, service.MediaContentUpdateInput{
		Title:       req.Title,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		Status:      req.Status,
		Views:       req.Views,
		Likes:       req.Likes,
		Comments:    req.Comments,
		Shares:      req.Shares,
		PublishTime: req.PublishTime,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, content)
}

func (h *MediaHandler) DeleteContent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_content_id", "Invalid content ID")
		return
	}

	err = h.mediaContentService.DeleteContent(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Content deleted successfully"})
}

func (h *MediaHandler) ListVersions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_content_id", "Invalid content ID")
		return
	}

	page, pageSize := response.ParsePagination(c)

	versions, total, err := h.mediaContentService.ListVersions(c.Request.Context(), id, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, versions, total, page, pageSize)
}

func (h *MediaHandler) GetVersion(c *gin.Context) {
	versionId, err := strconv.Atoi(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid_version_id", "Invalid version ID")
		return
	}

	version, err := h.mediaContentService.GetVersion(c.Request.Context(), versionId)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, version)
}

func (h *MediaHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	accounts := router.Group("/media/accounts")
	accounts.Use(authMiddleware)
	{
		accounts.GET("", h.ListAccounts)
		accounts.GET("/:id", h.GetAccount)
		accounts.POST("", h.CreateAccount)
		accounts.PUT("/:id", h.UpdateAccount)
		accounts.DELETE("/:id", h.DeleteAccount)
	}

	contents := router.Group("/media/contents")
	contents.Use(authMiddleware)
	{
		contents.GET("", h.ListContents)
		contents.GET("/:id", h.GetContent)
		contents.POST("", h.CreateContent)
		contents.PUT("/:id", h.UpdateContent)
		contents.DELETE("/:id", h.DeleteContent)
		contents.GET("/:id/versions", h.ListVersions)
		contents.GET("/versions/:versionId", h.GetVersion)
	}
}
