package phrases

type Language string

const (
	Fa Language = "fa"
	En Language = "en"
)

var errorMessagePhrase = map[Language]map[MessagePhrase]string{
	Fa: {
		// Default messages
		DefaultValidationID:       "درخواست نامعتبر است",
		DefaultNotFoundID:         "منبع یافت نشد",
		DefaultConflictID:         "درخواست با وضعیت فعلی در تضاد است",
		DefaultUnauthorizedID:     "احراز هویت نشده است",
		DefaultForbiddenID:        "دسترسی مجاز نیست",
		DefaultTimeoutID:          "درخواست منقضی شده است",
		DefaultRateLimitID:        "تعداد درخواست بیش از حد مجاز است",
		DefaultTooLargeID:         "حجم درخواست بیش از حد مجاز است",
		DefaultInternalID:         "خطای داخلی سرور",
		DefaultMethodNotAllowedID: "متد HTTP مجاز نیست",

		// User messages
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
		// Default messages
		DefaultValidationID:       "Invalid request",
		DefaultNotFoundID:         "Resource not found",
		DefaultConflictID:         "Request conflicts with current state",
		DefaultUnauthorizedID:     "Unauthorized",
		DefaultForbiddenID:        "Forbidden",
		DefaultTimeoutID:          "Request timeout",
		DefaultRateLimitID:        "Too many requests",
		DefaultTooLargeID:         "Request entity too large",
		DefaultInternalID:         "Internal server error",
		DefaultMethodNotAllowedID: "Method not allowed",

		// User messages
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

	// First, try to find message in requested language
	if msg, ok := errorMessagePhrase[lan][phrase]; ok {
		return msg
	}

	// If not found in requested language, search in all other languages
	// and return the first match found
	for lang, messages := range errorMessagePhrase {
		if lang != lan { // Skip the requested language (already checked)
			if msg, ok := messages[phrase]; ok {
				return msg
			}
		}
	}

	// If not found in any language, return default message
	return "Unknown error"
}
