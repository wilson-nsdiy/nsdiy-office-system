package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oa-nsdiy/backend/internal/middleware"
	"oa-nsdiy/backend/internal/pkg/response"
	"oa-nsdiy/backend/internal/service"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

type FileCreateRequest struct {
	Filename         string `json:"filename" binding:"required"`
	OriginalFilename string `json:"originalFilename" binding:"required"`
	FilePath         string `json:"filePath" binding:"required"`
	FileSize         int64  `json:"fileSize" binding:"required"`
	MimeType         string `json:"mimeType" binding:"required"`
	FileType         string `json:"fileType" binding:"required"`
	Extension        string `json:"extension" binding:"required"`
	Purpose          string `json:"purpose"`
	Md5              string `json:"md5"`
}

type FileUpdateRequest struct {
	Filename         string `json:"filename"`
	OriginalFilename string `json:"originalFilename"`
	FilePath         string `json:"filePath"`
	FileSize         int64  `json:"fileSize"`
	MimeType         string `json:"mimeType"`
	FileType         string `json:"fileType"`
	Extension        string `json:"extension"`
	Purpose          string `json:"purpose"`
	Md5              string `json:"md5"`
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	keyword := c.Query("keyword")
	fileType := c.Query("fileType")

	var uploaderID *int
	if uid := c.Query("uploaderId"); uid != "" {
		id, _ := strconv.Atoi(uid)
		uploaderID = &id
	}

	page, pageSize := response.ParsePagination(c)

	result, err := h.fileService.ListFiles(c.Request.Context(), keyword, fileType, uploaderID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, result.Items, result.Total, page, pageSize)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_file_id", "Invalid file ID")
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, file)
}

func (h *FileHandler) CreateFile(c *gin.Context) {
	var req FileCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	file, err := h.fileService.CreateFile(c.Request.Context(), service.FileCreateInput{
		Filename:         req.Filename,
		OriginalFilename: req.OriginalFilename,
		FilePath:         req.FilePath,
		FileSize:         req.FileSize,
		MimeType:         req.MimeType,
		FileType:         req.FileType,
		Extension:        req.Extension,
		Purpose:          req.Purpose,
		Md5:              req.Md5,
	}, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, file)
}

func (h *FileHandler) UpdateFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_file_id", "Invalid file ID")
		return
	}

	var req FileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid_request", "Invalid request body")
		return
	}

	file, err := h.fileService.UpdateFile(c.Request.Context(), id, service.FileUpdateInput{
		Filename:         req.Filename,
		OriginalFilename: req.OriginalFilename,
		FilePath:         req.FilePath,
		FileSize:         req.FileSize,
		MimeType:         req.MimeType,
		FileType:         req.FileType,
		Extension:        req.Extension,
		Purpose:          req.Purpose,
		Md5:              req.Md5,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, file)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid_file_id", "Invalid file ID")
		return
	}

	err = h.fileService.DeleteFile(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "File deleted successfully"})
}

func (h *FileHandler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	files := router.Group("/files")
	files.Use(authMiddleware)
	{
		files.GET("", h.ListFiles)
		files.GET("/:id", h.GetFile)
		files.POST("", h.CreateFile)
		files.PUT("/:id", h.UpdateFile)
		files.DELETE("/:id", h.DeleteFile)
	}
}
