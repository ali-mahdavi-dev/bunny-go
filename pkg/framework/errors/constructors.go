package errors

import (
	"fmt"

	"shikposh-backend/pkg/framework/errors/phrases"
)

// New creates a new error with the given parameters
func New(id string, errType ErrorType, message, detail string) Error {
	return &AppError{
		IDField:      id,
		TypeField:    errType,
		MessageField: message,
		DetailField:  detail,
	}
}

// Validation creates a validation error
func Validation(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultValidationID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeValidation,
		MessageField: message,
		DetailField:  "",
	}
}

// NotFound creates a not found error
func NotFound(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultNotFoundID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeNotFound,
		MessageField: message,
		DetailField:  "",
	}
}

// Conflict creates a conflict error
func Conflict(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultConflictID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeConflict,
		MessageField: message,
		DetailField:  "",
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultUnauthorizedID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeUnauthorized,
		MessageField: message,
		DetailField:  "",
	}
}

// Forbidden creates a forbidden error
func Forbidden(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultForbiddenID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeForbidden,
		MessageField: message,
		DetailField:  "",
	}
}

// Timeout creates a timeout error
func Timeout(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultTimeoutID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeTimeout,
		MessageField: message,
		DetailField:  "",
	}
}

// RateLimit creates a rate limit error
func RateLimit(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultRateLimitID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeRateLimit,
		MessageField: message,
		DetailField:  "",
	}
}

// TooLarge creates a request entity too large error
func TooLarge(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultTooLargeID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeTooLarge,
		MessageField: message,
		DetailField:  "",
	}
}

// Internal creates an internal server error
func Internal(detail string) Error {
	return &AppError{
		IDField:      string(phrases.DefaultInternalID),
		TypeField:    ErrorTypeInternal,
		MessageField: phrases.GetMessage(phrases.DefaultInternalID, ""),
		DetailField:  detail,
	}
}

// MethodNotAllowed creates a method not allowed error
func MethodNotAllowed(id phrases.MessagePhrase, args ...interface{}) Error {
	var message string
	if id == "" {
		id = phrases.DefaultMethodNotAllowedID
	}
	message = phrases.GetMessage(id, "")
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return &AppError{
		IDField:      string(id),
		TypeField:    ErrorTypeMethodNotAllowed,
		MessageField: message,
		DetailField:  "",
	}
}

