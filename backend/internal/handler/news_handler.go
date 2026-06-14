package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type NewsHandler struct {
	newsGroupService *service.NewsGroupService
	newsService      *service.NewsService
}

func NewNewsHandler(newsGroupService *service.NewsGroupService, newsService *service.NewsService) *NewsHandler {
	return &NewsHandler{
		newsGroupService: newsGroupService,
		newsService:      newsService,
	}
}

type NewsGroupCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   *int   `json:"sortOrder"`
}

type NewsGroupUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   *int   `json:"sortOrder"`
}

type NewsCreateRequest struct {
	GroupID int    `json:"groupId" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type NewsUpdateRequest struct {
	GroupID *int   `json:"groupId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *NewsHandler) GetAllGroups(c *gin.Context) {
	groups, err := h.newsGroupService.GetAllGroups(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, groups)
}

func (h *NewsHandler) GetGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_group_id", "Invalid group ID")
		return
	}

	group, err := h.newsGroupService.GetGroup(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, group)
}

func (h *NewsHandler) CreateGroup(c *gin.Context) {
	var req NewsGroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	group, err := h.newsGroupService.CreateGroup(c.Request.Context(), service.NewsGroupCreateInput{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   sortOrder,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, group)
}

func (h *NewsHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_group_id", "Invalid group ID")
		return
	}

	var req NewsGroupUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	group, err := h.newsGroupService.UpdateGroup(c.Request.Context(), id, service.NewsGroupUpdateInput{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   sortOrder,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, group)
}

func (h *NewsHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_group_id", "Invalid group ID")
		return
	}

	err = h.newsGroupService.DeleteGroup(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Group deleted successfully"})
}

func (h *NewsHandler) GetNewsList(c *gin.Context) {
	var groupID *int
	if gid := c.Query("groupId"); gid != "" {
		id, _ := strconv.Atoi(gid)
		groupID = &id
	}

	keyword := c.Query("keyword")
	page, pageSize := response.ParsePagination(c)

	result, err := h.newsService.GetNewsList(c.Request.Context(), groupID, keyword, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_news_id", "Invalid news ID")
		return
	}

	news, err := h.newsService.GetNews(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, news)
}

func (h *NewsHandler) CreateNews(c *gin.Context) {
	var req NewsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	news, err := h.newsService.CreateNews(c.Request.Context(), service.NewsCreateInput{
		GroupID: req.GroupID,
		Title:   req.Title,
		Content: req.Content,
	}, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, news)
}

func (h *NewsHandler) UpdateNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_news_id", "Invalid news ID")
		return
	}

	var req NewsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	news, err := h.newsService.UpdateNews(c.Request.Context(), id, service.NewsUpdateInput{
		GroupID: req.GroupID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, news)
}

func (h *NewsHandler) DeleteNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_news_id", "Invalid news ID")
		return
	}

	err = h.newsService.DeleteNews(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "News deleted successfully"})
}

func (h *NewsHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	newsGroups := router.Group("/news-groups")
	newsGroups.Use(authMiddleware)
	{
		newsGroups.GET("", h.GetAllGroups)
		newsGroups.GET("/:id", h.GetGroup)
		newsGroups.POST("", h.CreateGroup)
		newsGroups.PUT("/:id", h.UpdateGroup)
		newsGroups.DELETE("/:id", h.DeleteGroup)
	}

	news := router.Group("/news")
	news.Use(authMiddleware)
	{
		news.GET("", h.GetNewsList)
		news.GET("/:id", h.GetNews)
		news.POST("", h.CreateNews)
		news.PUT("/:id", h.UpdateNews)
		news.DELETE("/:id", h.DeleteNews)
	}
}
