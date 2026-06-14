package response

import (
	"math"
	"net/http"

	infraerrors "oa-nsdiy/backend/internal/pkg/errors"
	"oa-nsdiy/backend/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response 标准API响应格式
type Response struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	Reason   string            `json:"reason,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Data     any               `json:"data,omitempty"`
}

// PaginatedResponse 分页数据格式
type PaginatedResponse struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

// Success 返回成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created 返回创建成功响应
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithDetails returns an error response with reason and metadata fields.
func ErrorWithDetails(c *gin.Context, statusCode int, message, reason string, metadata map[string]string) {
	c.JSON(statusCode, Response{
		Code:     statusCode,
		Message:  message,
		Reason:   reason,
		Metadata: metadata,
	})
}

// ErrorFrom converts an ApplicationError (or any error) into the error response.
// It returns true if an error was written.
func ErrorFrom(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	statusCode, status := infraerrors.ToHTTP(err)

	// Log internal errors with full details for debugging
	if statusCode >= 500 && c.Request != nil {
		logger.FromContext(c.Request.Context()).Error("Internal server error",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Error(err),
		)
	}

	ErrorWithDetails(c, statusCode, status.Message, status.Reason, status.Metadata)
	return true
}

// BadRequest 返回400错误，reason 为业务语义错误码，message 为用户可读描述
func BadRequest(c *gin.Context, reason, message string) {
	ErrorFrom(c, infraerrors.BadRequest(reason, message))
}

// Unauthorized 返回401错误
func Unauthorized(c *gin.Context, reason, message string) {
	ErrorFrom(c, infraerrors.Unauthorized(reason, message))
}

// Forbidden 返回403错误
func Forbidden(c *gin.Context, reason, message string) {
	ErrorFrom(c, infraerrors.Forbidden(reason, message))
}

// NotFound 返回404错误
func NotFound(c *gin.Context, reason, message string) {
	ErrorFrom(c, infraerrors.NotFound(reason, message))
}

// Conflict 返回409错误
func Conflict(c *gin.Context, reason, message string) {
	ErrorFrom(c, infraerrors.Conflict(reason, message))
}

// InternalError 返回500错误
func InternalError(c *gin.Context, reason, message string) {
	ErrorFrom(c, infraerrors.Internal(reason, message))
}

// Paginated 返回分页数据
func Paginated(c *gin.Context, items any, total int64, page, pageSize int) {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	Success(c, PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// ParsePagination 解析分页参数，返回 page 和 pageSize
func ParsePagination(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 10

	if p := c.Query("page"); p != "" {
		if val, err := parseInt(p); err == nil && val > 0 {
			page = val
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if val, err := parseInt(ps); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	} else if l := c.Query("limit"); l != "" {
		if val, err := parseInt(l); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

	return page, pageSize
}

func parseInt(s string) (int, error) {
	var result int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		result = result*10 + int(c-'0')
	}
	return result, nil
}
