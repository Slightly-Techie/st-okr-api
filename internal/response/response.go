package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents standard error codes
type ErrorCode string

const (
	// Client Error Codes (4xx)
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeConflict         ErrorCode = "CONFLICT"
	ErrCodeBadRequest       ErrorCode = "BAD_REQUEST"
	ErrCodeTooManyRequests  ErrorCode = "TOO_MANY_REQUESTS"

	// Server Error Codes (5xx)
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeDatabaseError      ErrorCode = "DATABASE_ERROR"
)

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	RequestID string      `json:"request_id"`
}

// SuccessResponse represents the success response structure
type SuccessResponse struct {
	Success   bool   `json:"success"`
	Data      any    `json:"data,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id"`
	Meta      *Meta  `json:"meta,omitempty"`
}

// Meta represents metadata for responses (pagination, timestamps, etc.)
type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	Page      *int      `json:"page,omitempty"`
	Limit     *int      `json:"limit,omitempty"`
	Total     *int64    `json:"total,omitempty"`
	HasMore   *bool     `json:"has_more,omitempty"`
}

// getRequestID extracts request ID from context
func getRequestID(c *gin.Context) string {
	if reqID, exists := c.Get("request_id"); exists {
		return reqID.(string)
	}
	return "unknown"
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, data any, message string) {
	response := SuccessResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		RequestID: getRequestID(c),
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
	}
	c.JSON(statusCode, response)
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(c *gin.Context, statusCode int, data any, message string, meta *Meta) {
	if meta == nil {
		meta = &Meta{}
	}
	meta.Timestamp = time.Now().UTC()

	response := SuccessResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		RequestID: getRequestID(c),
		Meta:      meta,
	}
	c.JSON(statusCode, response)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code ErrorCode, message string, details map[string]string) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: getRequestID(c),
	}
	c.JSON(statusCode, response)
}

// BadRequest sends a 400 bad request error
func BadRequest(c *gin.Context, message string, details map[string]string) {
	Error(c, http.StatusBadRequest, ErrCodeBadRequest, message, details)
}

// ValidationError sends a 400 validation error
func ValidationError(c *gin.Context, message string, details map[string]string) {
	Error(c, http.StatusBadRequest, ErrCodeValidationFailed, message, details)
}

// Unauthorized sends a 401 unauthorized error
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrCodeUnauthorized, message, nil)
}

// Forbidden sends a 403 forbidden error
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, ErrCodeForbidden, message, nil)
}

// NotFound sends a 404 not found error
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, ErrCodeNotFound, message, nil)
}

// Conflict sends a 409 conflict error
func Conflict(c *gin.Context, message string, details map[string]string) {
	Error(c, http.StatusConflict, ErrCodeConflict, message, details)
}

// InternalError sends a 500 internal server error
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeInternalError, message, nil)
}

// DatabaseError sends a 500 database error
func DatabaseError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeDatabaseError, message, nil)
}

// Created sends a 201 created response
func Created(c *gin.Context, data any, message string) {
	Success(c, http.StatusCreated, data, message)
}

// OK sends a 200 success response
func OK(c *gin.Context, data any, message string) {
	Success(c, http.StatusOK, data, message)
}

// NoContent sends a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
