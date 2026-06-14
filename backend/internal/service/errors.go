package service

import (
	"net/http"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/internal/pkg/errors"
)

// ServiceError is an alias for ApplicationError.
// Service 层直接返回 *ApplicationError，Handler 用 response.ErrorFrom 统一渲染。
type ServiceError = errors.ApplicationError

// BadRequestErr 创建 400 错误。
func BadRequestErr(reason, message string) *ServiceError {
	return errors.BadRequest(reason, message)
}

// UnauthorizedErr 创建 401 错误。
func UnauthorizedErr(reason, message string) *ServiceError {
	return errors.Unauthorized(reason, message)
}

// ForbiddenErr 创建 403 错误。
func ForbiddenErr(reason, message string) *ServiceError {
	return errors.Forbidden(reason, message)
}

// NotFoundErr 创建 404 错误。
func NotFoundErr(reason, message string) *ServiceError {
	return errors.NotFound(reason, message)
}

// ConflictErr 创建 409 错误。
func ConflictErr(reason, message string) *ServiceError {
	return errors.Conflict(reason, message)
}

// InternalErr 创建 500 错误。
func InternalErr(reason, message string) *ServiceError {
	return errors.Internal(reason, message)
}

// NewServiceError 为兼容保留的通用构造函数。
func NewServiceError(code int, reason, message string) *ServiceError {
	return errors.New(code, reason, message)
}

// IsStatus 检查 err 的 HTTP 状态码是否匹配。
func IsStatus(err error, code int) bool {
	if err == nil {
		return false
	}
	return errors.Code(err) == code
}

// IsNotFound 快速判断是否为 404。
func IsNotFound(err error) bool {
	return IsStatus(err, http.StatusNotFound)
}

// IsBadRequest 快速判断是否为 400。
func IsBadRequest(err error) bool {
	return IsStatus(err, http.StatusBadRequest)
}

// HandleRepoErr converts a repository error to the appropriate ServiceError.
// Returns NotFoundErr for ent not-found errors, otherwise passes through.
func HandleRepoErr(err error, reason, message string) error {
	if err == nil {
		return nil
	}
	if ent.IsNotFound(err) {
		return NotFoundErr(reason, message)
	}
	// If it's already an ApplicationError, return as-is
	if _, ok := err.(*errors.ApplicationError); ok {
		return err
	}
	return err
}
