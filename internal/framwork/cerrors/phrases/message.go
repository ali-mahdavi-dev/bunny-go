package phrases

type MessagePhrase string

var (
	UserNotFound      MessagePhrase = "User.NotFound"
	UserAlreadyExists MessagePhrase = "User.AlreadyExists"
	UserAgeInvalid    MessagePhrase = "User.AgeInvalid"
	UserInvalid       MessagePhrase = "User.Invalid"
	OperationCanNot   MessagePhrase = "Operation.CanNot"
	FailedParseJson   MessagePhrase = "FailedParseJson"
	FailedParseQuery  MessagePhrase = "FailedParseQuery"
	FailedParseForm   MessagePhrase = "FailedParseForm"
)

type Language string

const (
	Fa Language = "fa"
	En Language = "en"
)

var errorMessagePhrase = map[Language]map[MessagePhrase]string{
	Fa: {
		UserNotFound:      "کاربر پیدا نشد",
		UserAlreadyExists: "کاربر از قبل وجود دارد",
		UserAgeInvalid:    "سن کاربر کمتر از ۱۸ است",
		UserInvalid:       "اطلاعات کاربر درست نمیباشد",
		OperationCanNot:   "عملیات موفق آمیز نبود. لطفا دوباره تلاش بفرمایید",
		FailedParseJson:   "خطا در تجزیه JSON: %s",
		FailedParseQuery:  "خطا در تجزیه Query: %s",
		FailedParseForm:   "خطا در تجزیه Form: %s",
	},
	En: {
		UserNotFound:      "User not found",
		UserAlreadyExists: "User already exists",
		UserAgeInvalid:    "User age is less than 18",
		UserInvalid:       "User information is not valid",
		OperationCanNot:   "Operation was not successful. Please try again",
		FailedParseJson:   "Failed to parse json: %s",
		FailedParseQuery:  "Failed to parse query: %s",
		FailedParseForm:   "Failed to parse form: %s",
	},
}

func GetMessage(phrase MessagePhrase, lan Language) string {
	if lan == "" {
		lan = Fa
	}

	if msg, ok := errorMessagePhrase[lan][phrase]; ok {
		return msg
	}
	return "Unknown error"
}
