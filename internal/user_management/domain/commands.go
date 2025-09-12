package domain

// user
type RegisterUser struct {
	AvatarIdentifier string `json:"avatar_identifier"`
	UserName         string `json:"user_name"`
	Password         string `json:password`
}
