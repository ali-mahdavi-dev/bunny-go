package phrases

type MessagePhrase string

const (
	// Default error IDs
	DefaultValidationID       MessagePhrase = "validation"
	DefaultNotFoundID         MessagePhrase = "not_found"
	DefaultConflictID         MessagePhrase = "conflict"
	DefaultUnauthorizedID     MessagePhrase = "unauthorized"
	DefaultForbiddenID        MessagePhrase = "forbidden"
	DefaultTimeoutID          MessagePhrase = "timeout"
	DefaultRateLimitID        MessagePhrase = "rate_limit"
	DefaultTooLargeID         MessagePhrase = "too_large"
	DefaultInternalID         MessagePhrase = "internal_server_error"
	DefaultMethodNotAllowedID MessagePhrase = "method_not_allowed"
)

var (
	// User
	UserNotFound      MessagePhrase = "User.NotFound"
	UserAlreadyExists MessagePhrase = "User.AlreadyExists"
	UserAgeInvalid    MessagePhrase = "User.AgeInvalid"
	UserInvalid       MessagePhrase = "User.Invalid"

	// Operation
	OperationCanNot MessagePhrase = "Operation.CanNot"

	// Failed
	FailedParseJson  MessagePhrase = "FailedParseJson"
	FailedParseQuery MessagePhrase = "FailedParseQuery"
	FailedParseForm  MessagePhrase = "FailedParseForm"
)
