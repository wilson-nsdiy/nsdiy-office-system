package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

type ArticleCreateRequest struct {
	Title            string `json:"title" binding:"required"`
	Content          string `json:"content"`
	Summary          string `json:"summary"`
	Status           string `json:"status"`
	CoverDescription string `json:"coverDescription"`
	CoverUrl         string `json:"coverUrl"`
}

type ArticleUpdateRequest struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	Summary          string `json:"summary"`
	Status           string `json:"status"`
	CoverDescription string `json:"coverDescription"`
	CoverUrl         string `json:"coverUrl"`
	EditReason       string `json:"editReason"`
}

func (h *ArticleHandler) GetArticleList(c *gin.Context) {
	keyword := c.Query("keyword")
	status := c.Query("status")
	page, pageSize := response.ParsePagination(c)

	result, err := h.articleService.GetArticleList(c.Request.Context(), keyword, status, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_article_id", "Invalid article ID")
		return
	}

	article, err := h.articleService.GetArticle(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, article)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req ArticleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	article, err := h.articleService.CreateArticle(c.Request.Context(), service.ArticleCreateInput{
		Title:            req.Title,
		Content:          req.Content,
		Summary:          req.Summary,
		Status:           req.Status,
		CoverDescription: req.CoverDescription,
		CoverUrl:         req.CoverUrl,
	}, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, article)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_article_id", "Invalid article ID")
		return
	}

	var req ArticleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	article, err := h.articleService.UpdateArticle(c.Request.Context(), id, service.ArticleUpdateInput{
		Title:            req.Title,
		Content:          req.Content,
		Summary:          req.Summary,
		Status:           req.Status,
		CoverDescription: req.CoverDescription,
		CoverUrl:         req.CoverUrl,
		EditReason:       req.EditReason,
	}, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_article_id", "Invalid article ID")
		return
	}

	err = h.articleService.DeleteArticle(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Article deleted successfully"})
}

func (h *ArticleHandler) GetArticleVersions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_article_id", "Invalid article ID")
		return
	}

	page, pageSize := response.ParsePagination(c)

	versions, total, err := h.articleService.GetArticleVersions(c.Request.Context(), id, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, versions, total, page, pageSize)
}

func (h *ArticleHandler) GetArticleVersion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_article_id", "Invalid article ID")
		return
	}

	versionId, err := strconv.Atoi(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid_version_id", "Invalid version ID")
		return
	}

	version, err := h.articleService.GetArticleVersion(c.Request.Context(), id, versionId)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, version)
}

func (h *ArticleHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	articles := router.Group("/articles")
	articles.Use(authMiddleware)
	{
		articles.GET("", h.GetArticleList)
		articles.GET("/:id", h.GetArticle)
		articles.POST("", h.CreateArticle)
		articles.PUT("/:id", h.UpdateArticle)
		articles.DELETE("/:id", h.DeleteArticle)
		articles.GET("/:id/versions", h.GetArticleVersions)
		articles.GET("/:id/versions/:versionId", h.GetArticleVersion)
	}
}
