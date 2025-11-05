package errors

// ErrorType represents the type/category of an error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors (invalid input)
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound represents resource not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeConflict represents conflict errors (e.g., duplicate resource)
	ErrorTypeConflict ErrorType = "conflict"
	// ErrorTypeUnauthorized represents authentication errors
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	// ErrorTypeForbidden represents authorization errors
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeTimeout represents timeout errors
	ErrorTypeTimeout ErrorType = "timeout"
	// ErrorTypeRateLimit represents rate limiting errors
	ErrorTypeRateLimit ErrorType = "rate_limit"
	// ErrorTypeTooLarge represents request entity too large errors
	ErrorTypeTooLarge ErrorType = "too_large"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeMethodNotAllowed represents HTTP method not allowed errors
	ErrorTypeMethodNotAllowed ErrorType = "method_not_allowed"
)

// Error is the core error interface that is independent of HTTP
type Error interface {
	// ID returns a unique identifier for this error type
	ID() string
	// Type returns the error type/category
	Type() ErrorType
	// Message returns the user-friendly error message
	Message() string
	// Detail returns additional error details (optional)
	Detail() string
	// Error returns the string representation (implements error interface)
	Error() string
}

// AppError is the base implementation of Error
type AppError struct {
	IDField      string
	TypeField    ErrorType
	MessageField string
	DetailField  string
}

func (e *AppError) ID() string      { return e.IDField }
func (e *AppError) Type() ErrorType { return e.TypeField }
func (e *AppError) Message() string { return e.MessageField }
func (e *AppError) Detail() string  { return e.DetailField }
func (e *AppError) Error() string {
	if e.DetailField != "" {
		return e.MessageField + ": " + e.DetailField
	}
	return e.MessageField
}
