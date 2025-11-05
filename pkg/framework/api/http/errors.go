package http

import (
	"encoding/json"
	"net/http"

	"shikposh-backend/pkg/framework/errors"
	"shikposh-backend/pkg/framework/errors/phrases"
)

// HTTPError is the HTTP-specific error representation
type HTTPError struct {
	Code    string `json:"code"`    // Error code (from error ID)
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Status  string `json:"status"`  // HTTP status text
}

// ErrorToHTTP converts an app error to HTTP error with appropriate status code
func ErrorToHTTP(err errors.Error) Error {
	statusCode := errorTypeToHTTPStatus(err.Type())

	httpErr := HTTPError{
		Code:    err.ID(),           // Use error ID as code
		Message: err.Message(),
		Detail:  err.Detail(),
		Status:  http.StatusText(statusCode),
	}
	
	return &httpErrorAdapter{httpError: httpErr}
}

// errorTypeToHTTPStatus maps error types to HTTP status codes
func errorTypeToHTTPStatus(errType errors.ErrorType) int {
	switch errType {
	case errors.ErrorTypeValidation:
		return http.StatusBadRequest
	case errors.ErrorTypeNotFound:
		return http.StatusNotFound
	case errors.ErrorTypeConflict:
		return http.StatusConflict
	case errors.ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrorTypeForbidden:
		return http.StatusForbidden
	case errors.ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case errors.ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case errors.ErrorTypeTooLarge:
		return http.StatusRequestEntityTooLarge
	case errors.ErrorTypeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case errors.ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// HTTPError interface for backward compatibility
type Error interface {
	Code() string
	Message() string
	Detail() string
	Status() string
	Error() string
}

// httpErrorAdapter adapts HTTPError to Error interface
type httpErrorAdapter struct {
	httpError HTTPError
}

func (e *httpErrorAdapter) Code() string     { return e.httpError.Code }
func (e *httpErrorAdapter) Message() string  { return e.httpError.Message }
func (e *httpErrorAdapter) Detail() string   { return e.httpError.Detail }
func (e *httpErrorAdapter) Status() string   { return e.httpError.Status }
func (e *httpErrorAdapter) Error() string {
	b, _ := json.Marshal(e.httpError)
	return string(b)
}

// ToHTTPError converts an Error interface to HTTPError struct
func ToHTTPError(err Error) *HTTPError {
	return &HTTPError{
		Code:    err.Code(),
		Message: err.Message(),
		Detail:  err.Detail(),
		Status:  err.Status(),
	}
}

// Convenience functions that wrap app errors and convert to HTTP errors

// BadRequest creates a validation error and converts to HTTP
func BadRequest(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.Validation(id, args...)
	return ErrorToHTTP(appErr)
}

// NotFound creates a not found error and converts to HTTP
func NotFound(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.NotFound(id, args...)
	return ErrorToHTTP(appErr)
}

// Conflict creates a conflict error and converts to HTTP
func Conflict(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.Conflict(id, args...)
	return ErrorToHTTP(appErr)
}

// Unauthorized creates an unauthorized error and converts to HTTP
func Unauthorized(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.Unauthorized(id, args...)
	return ErrorToHTTP(appErr)
}

// Forbidden creates a forbidden error and converts to HTTP
func Forbidden(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.Forbidden(id, args...)
	return ErrorToHTTP(appErr)
}

// Timeout creates a timeout error and converts to HTTP
func Timeout(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.Timeout(id, args...)
	return ErrorToHTTP(appErr)
}

// TooManyRequests creates a rate limit error and converts to HTTP
func TooManyRequests(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.RateLimit(id, args...)
	return ErrorToHTTP(appErr)
}

// RequestEntityTooLarge creates a too large error and converts to HTTP
func RequestEntityTooLarge(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.TooLarge(id, args...)
	return ErrorToHTTP(appErr)
}

// MethodNotAllowed creates a method not allowed error and converts to HTTP
func MethodNotAllowed(id phrases.MessagePhrase, args ...interface{}) Error {
	appErr := errors.MethodNotAllowed(id, args...)
	return ErrorToHTTP(appErr)
}

// InternalServerError creates an internal error and converts to HTTP
func InternalServerError(detail string) Error {
	appErr := errors.Internal(detail)
	return ErrorToHTTP(appErr)
}

