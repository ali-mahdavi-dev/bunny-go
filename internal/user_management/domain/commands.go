package domain

// user
type RegisterUser struct {
	AvatarIdentifier string `json:"avatar_identifier"`
	UserName         string `json:"username"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
}
