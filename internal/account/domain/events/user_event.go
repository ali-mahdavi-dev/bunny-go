package events

// user
type RegisterUserEvent struct {
	AvatarIdentifier string `json:"avatar_identifier" binding:"required"`
	UserName         string `json:"user_name" binding:"required,min=3"`
	FirstName        string `json:"first_name" binding:"required,min=3"`
	LastName         string `json:"last_name" binding:"required,min=3"`
	Email            string `json:"email" binding:"required,min=3"`
	Password         string `json:"password" binding:"required,min=3"`
}
