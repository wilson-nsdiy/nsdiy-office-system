package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type TaskHandler struct {
	taskService    *service.TaskService
	projectService *service.ProjectService
}

func NewTaskHandler(taskService *service.TaskService, projectService *service.ProjectService) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		projectService: projectService,
	}
}

type TaskCreateRequest struct {
	Title            string  `json:"title" binding:"required"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
	Priority         string  `json:"priority"`
	AssigneeID       *int    `json:"assigneeId"`
	ParentID         *int    `json:"parentId"`
	PlannedStartDate string  `json:"plannedStartDate"`
	PlannedEndDate   string  `json:"plannedEndDate"`
	EstimatedHours   float64 `json:"estimatedHours"`
}

type TaskUpdateRequest struct {
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
	Priority         string  `json:"priority"`
	AssigneeID       *int    `json:"assigneeId"`
	ParentID         *int    `json:"parentId"`
	PlannedStartDate string  `json:"plannedStartDate"`
	PlannedEndDate   string  `json:"plannedEndDate"`
	EstimatedHours   float64 `json:"estimatedHours"`
}

func (h *TaskHandler) GetTaskList(c *gin.Context) {
	var projectID *int
	if pid := c.Query("projectId"); pid != "" {
		id, _ := strconv.Atoi(pid)
		projectID = &id
	}

	status := c.Query("status")
	priority := c.Query("priority")

	var assigneeID *int
	if aid := c.Query("assigneeId"); aid != "" {
		id, _ := strconv.Atoi(aid)
		assigneeID = &id
	}

	page, pageSize := response.ParsePagination(c)

	result, err := h.taskService.GetTaskList(c.Request.Context(), projectID, status, priority, assigneeID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_task_id", "Invalid task ID")
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil {
		response.BadRequest(c, "invalid_project_id", "Invalid project ID")
		return
	}

	// Check project membership
	userID := middleware.GetUserID(c)
	isMember, err := h.projectService.IsMember(c.Request.Context(), "", projectId) // TODO: pass projectNo
	if err != nil || !isMember {
		response.Forbidden(c, "not_project_member", "Not a project member")
		return
	}

	var req TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), service.TaskCreateInput{
		Title:            req.Title,
		Description:      req.Description,
		Status:           req.Status,
		Priority:         req.Priority,
		AssigneeID:       req.AssigneeID,
		ParentID:         req.ParentID,
		PlannedStartDate: req.PlannedStartDate,
		PlannedEndDate:   req.PlannedEndDate,
		EstimatedHours:   req.EstimatedHours,
	}, projectId, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_task_id", "Invalid task ID")
		return
	}

	var req TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), id, service.TaskUpdateInput{
		Title:            req.Title,
		Description:      req.Description,
		Status:           req.Status,
		Priority:         req.Priority,
		AssigneeID:       req.AssigneeID,
		ParentID:         req.ParentID,
		PlannedStartDate: req.PlannedStartDate,
		PlannedEndDate:   req.PlannedEndDate,
		EstimatedHours:   req.EstimatedHours,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_task_id", "Invalid task ID")
		return
	}

	err = h.taskService.DeleteTask(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	tasks := router.Group("/tasks")
	tasks.Use(authMiddleware)
	{
		tasks.GET("", h.GetTaskList)
		tasks.GET("/:id", h.GetTask)
		tasks.POST("", h.CreateTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}
