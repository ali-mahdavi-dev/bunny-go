package phrases

type MessagePhrase string

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
