package cerrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/cerrors/phrases"

	"github.com/pkg/errors"
)

// Define alias
var (
	WithStack = errors.WithStack
	Wrap      = errors.Wrap
	Wrapf     = errors.Wrapf
	Is        = errors.Is
	Errorf    = errors.Errorf
)

const (
	DefaultBadRequestID            = "bad_request"
	DefaultUnauthorizedID          = "unauthorized"
	DefaultForbiddenID             = "forbidden"
	DefaultNotFoundID              = "not_found"
	DefaultMethodNotAllowedID      = "method_not_allowed"
	DefaultTooManyRequestsID       = "too_many_requests"
	DefaultRequestEntityTooLargeID = "request_entity_too_large"
	DefaultInternalServerErrorID   = "internal_server_error"
	DefaultConflictID              = "conflict"
	DefaultRequestTimeoutID        = "request_timeout"
)

// Customize the error structure for implementation errors.Error interface
type Error interface {
	ID() string
	Code() int32
	Message() string
	Detail() string
	Status() string
	Error() string
}

type customeError struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *customeError) ID() string      { return e.IDField }
func (e *customeError) Code() int32     { return e.CodeField }
func (e *customeError) Message() string { return e.MessageField }
func (e *customeError) Detail() string  { return e.DetailField }
func (e *customeError) Status() string  { return e.StatusField }

func (e *customeError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func New(id string, code int32, message, detail, status string) Error {
	return &customeError{
		IDField:      id,
		CodeField:    code,
		MessageField: message,
		DetailField:  detail,
		StatusField:  status,
	}
}

// BadRequest generates a 400 error.
func BadRequest(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultBadRequestID
	}
	message = phrases.GetMessage(code, "")

	return &BadRequestErr{
		IDField:      string(code),
		CodeField:    http.StatusBadRequest,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusBadRequest),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultUnauthorizedID
	}
	message = phrases.GetMessage(code, "")
	return &UnauthorizedErr{
		IDField:      string(code),
		CodeField:    http.StatusUnauthorized,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusUnauthorized),
	}
}

// Forbidden generates a 403 error.
func Forbidden(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultForbiddenID
	}
	message = phrases.GetMessage(code, "")
	return &ForbiddenErr{
		IDField:      string(code),
		CodeField:    http.StatusForbidden,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusForbidden),
	}
}

// NotFound generates a 404 error.
func NotFound(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultNotFoundID
	}
	message = phrases.GetMessage(code, "")
	return &NotFoundErr{
		IDField:      string(code),
		CodeField:    http.StatusNotFound,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusNotFound),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultMethodNotAllowedID
	}
	message = phrases.GetMessage(code, "")
	return &MethodNotAllowedErr{
		IDField:      string(code),
		CodeField:    http.StatusMethodNotAllowed,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusMethodNotAllowed),
	}
}

// TooManyRequests generates a 429 error.
func TooManyRequests(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultTooManyRequestsID
	}
	message = phrases.GetMessage(code, "")
	return &TooManyRequestsErr{
		IDField:      string(code),
		CodeField:    http.StatusTooManyRequests,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusTooManyRequests),
	}
}

// Timeout generates a 408 error.
func Timeout(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultRequestTimeoutID
	}
	message = phrases.GetMessage(code, "")
	return &TimeoutErr{
		IDField:      string(code),
		CodeField:    http.StatusRequestTimeout,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusRequestTimeout),
	}
}

// Conflict generates a 409 error.
func Conflict(code, format string, a ...interface{}) Error {
	if code == "" {
		code = DefaultConflictID
	}
	return &ConflictErr{
		IDField:      string(code),
		CodeField:    http.StatusConflict,
		MessageField: fmt.Sprintf(format, a...),
		StatusField:  http.StatusText(http.StatusConflict),
	}
}

// RequestEntityTooLarge generates a 413 error.
func RequestEntityTooLarge(code phrases.MessagePhrase, a ...interface{}) Error {
	var message string
	if code == "" {
		code = DefaultRequestEntityTooLargeID
	}
	message = phrases.GetMessage(code, "")
	return &RequestEntityTooLargeErr{
		IDField:      string(code),
		CodeField:    http.StatusRequestEntityTooLarge,
		MessageField: fmt.Sprintf(message, a...),
		StatusField:  http.StatusText(http.StatusRequestEntityTooLarge),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(detail string) Error {
	return &InternalServerErr{
		IDField:     DefaultInternalServerErrorID,
		CodeField:   http.StatusInternalServerError,
		DetailField: detail,
		StatusField: http.StatusText(http.StatusInternalServerError),
	}
}

// As finds the first error in err's chain that matches *Error
func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var merr *Error
	if errors.As(err, &merr) {
		return merr, true
	}
	return nil, false
}

type BadRequestErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *BadRequestErr) ID() string      { return e.IDField }
func (e *BadRequestErr) Code() int32     { return e.CodeField }
func (e *BadRequestErr) Message() string { return e.MessageField }
func (e *BadRequestErr) Detail() string  { return e.DetailField }
func (e *BadRequestErr) Status() string  { return e.StatusField }
func (e *BadRequestErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type UnauthorizedErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *UnauthorizedErr) ID() string      { return e.IDField }
func (e *UnauthorizedErr) Code() int32     { return e.CodeField }
func (e *UnauthorizedErr) Message() string { return e.MessageField }
func (e *UnauthorizedErr) Detail() string  { return e.DetailField }
func (e *UnauthorizedErr) Status() string  { return e.StatusField }
func (e *UnauthorizedErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type ForbiddenErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *ForbiddenErr) ID() string      { return e.IDField }
func (e *ForbiddenErr) Code() int32     { return e.CodeField }
func (e *ForbiddenErr) Message() string { return e.MessageField }
func (e *ForbiddenErr) Detail() string  { return e.DetailField }
func (e *ForbiddenErr) Status() string  { return e.StatusField }
func (e *ForbiddenErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type NotFoundErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *NotFoundErr) ID() string      { return e.IDField }
func (e *NotFoundErr) Code() int32     { return e.CodeField }
func (e *NotFoundErr) Message() string { return e.MessageField }
func (e *NotFoundErr) Detail() string  { return e.DetailField }
func (e *NotFoundErr) Status() string  { return e.StatusField }
func (e *NotFoundErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type MethodNotAllowedErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *MethodNotAllowedErr) ID() string      { return e.IDField }
func (e *MethodNotAllowedErr) Code() int32     { return e.CodeField }
func (e *MethodNotAllowedErr) Message() string { return e.MessageField }
func (e *MethodNotAllowedErr) Detail() string  { return e.DetailField }
func (e *MethodNotAllowedErr) Status() string  { return e.StatusField }
func (e *MethodNotAllowedErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type TooManyRequestsErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *TooManyRequestsErr) ID() string      { return e.IDField }
func (e *TooManyRequestsErr) Code() int32     { return e.CodeField }
func (e *TooManyRequestsErr) Message() string { return e.MessageField }
func (e *TooManyRequestsErr) Detail() string  { return e.DetailField }
func (e *TooManyRequestsErr) Status() string  { return e.StatusField }
func (e *TooManyRequestsErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type RequestEntityTooLargeErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *RequestEntityTooLargeErr) ID() string      { return e.IDField }
func (e *RequestEntityTooLargeErr) Code() int32     { return e.CodeField }
func (e *RequestEntityTooLargeErr) Message() string { return e.MessageField }
func (e *RequestEntityTooLargeErr) Detail() string  { return e.DetailField }
func (e *RequestEntityTooLargeErr) Status() string  { return e.StatusField }
func (e *RequestEntityTooLargeErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type InternalServerErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *InternalServerErr) ID() string      { return e.IDField }
func (e *InternalServerErr) Code() int32     { return e.CodeField }
func (e *InternalServerErr) Message() string { return e.MessageField }
func (e *InternalServerErr) Detail() string  { return e.DetailField }
func (e *InternalServerErr) Status() string  { return e.StatusField }
func (e *InternalServerErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type ConflictErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *ConflictErr) ID() string      { return e.IDField }
func (e *ConflictErr) Code() int32     { return e.CodeField }
func (e *ConflictErr) Message() string { return e.MessageField }
func (e *ConflictErr) Detail() string  { return e.DetailField }
func (e *ConflictErr) Status() string  { return e.StatusField }
func (e *ConflictErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type TimeoutErr struct {
	IDField      string `json:"id"`
	CodeField    int32  `json:"code"`
	MessageField string `json:"message"`
	DetailField  string `json:"detail"`
	StatusField  string `json:"status"`
}

func (e *TimeoutErr) ID() string      { return e.IDField }
func (e *TimeoutErr) Code() int32     { return e.CodeField }
func (e *TimeoutErr) Message() string { return e.MessageField }
func (e *TimeoutErr) Detail() string  { return e.DetailField }
func (e *TimeoutErr) Status() string  { return e.StatusField }
func (e *TimeoutErr) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
