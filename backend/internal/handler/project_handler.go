package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

type ProjectCreateRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	Priority          string `json:"priority"`
	ExpectedStartDate string `json:"expectedStartDate"`
	ExpectedEndDate   string `json:"expectedEndDate"`
}

type ProjectUpdateRequest struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	Priority          string `json:"priority"`
	ExpectedStartDate string `json:"expectedStartDate"`
	ExpectedEndDate   string `json:"expectedEndDate"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
}

type ProjectMemberCreateRequest struct {
	UserId int    `json:"userId" binding:"required"`
	Role   string `json:"role"`
}

type ProjectMemberUpdateRequest struct {
	Role string `json:"role" binding:"required"`
}

func (h *ProjectHandler) GetProjectList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	keyword := c.Query("keyword")
	page, pageSize := response.ParsePagination(c)

	result, err := h.projectService.GetProjectList(c.Request.Context(), userID, keyword, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectNo := c.Param("projectNo")

	project, err := h.projectService.GetProject(c.Request.Context(), projectNo)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, project)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req ProjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	project, err := h.projectService.CreateProject(c.Request.Context(), service.ProjectCreateInput{
		Name:              req.Name,
		Description:       req.Description,
		Priority:          req.Priority,
		ExpectedStartDate: req.ExpectedStartDate,
		ExpectedEndDate:   req.ExpectedEndDate,
	}, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectNo := c.Param("projectNo")

	var req ProjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), projectNo, service.ProjectUpdateInput{
		Name:              req.Name,
		Description:       req.Description,
		Status:            req.Status,
		Priority:          req.Priority,
		ExpectedStartDate: req.ExpectedStartDate,
		ExpectedEndDate:   req.ExpectedEndDate,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectNo := c.Param("projectNo")

	err := h.projectService.DeleteProject(c.Request.Context(), projectNo)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Project deleted successfully"})
}

func (h *ProjectHandler) GetMembers(c *gin.Context) {
	projectNo := c.Param("projectNo")

	members, err := h.projectService.GetMembers(c.Request.Context(), projectNo)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, members)
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectNo := c.Param("projectNo")

	var req ProjectMemberCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	err := h.projectService.AddMember(c.Request.Context(), projectNo, req.UserId, req.Role)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, gin.H{"message": "Member added successfully"})
}

func (h *ProjectHandler) UpdateMemberRole(c *gin.Context) {
	projectNo := c.Param("projectNo")
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "invalid_user_id", "Invalid user ID")
		return
	}

	var req ProjectMemberUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	err = h.projectService.UpdateMemberRole(c.Request.Context(), projectNo, userId, req.Role)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Member role updated successfully"})
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectNo := c.Param("projectNo")
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "invalid_user_id", "Invalid user ID")
		return
	}

	err = h.projectService.RemoveMember(c.Request.Context(), projectNo, userId)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Member removed successfully"})
}

func (h *ProjectHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	projects := router.Group("/projects")
	projects.Use(authMiddleware)
	{
		projects.GET("", h.GetProjectList)
		projects.GET("/:projectNo", h.GetProject)
		projects.POST("", h.CreateProject)
		projects.PUT("/:projectNo", h.UpdateProject)
		projects.DELETE("/:projectNo", h.DeleteProject)
		projects.GET("/:projectNo/members", h.GetMembers)
		projects.POST("/:projectNo/members", h.AddMember)
		projects.PUT("/:projectNo/members/:userId", h.UpdateMemberRole)
		projects.DELETE("/:projectNo/members/:userId", h.RemoveMember)
	}
}
