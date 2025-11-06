package events

// user
type RegisterUserEvent struct {
	UserID           uint64 `json:"user_id"`
	AvatarIdentifier string `json:"avatar_identifier"`
	UserName         string `json:"user_name"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
}
