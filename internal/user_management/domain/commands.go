package domain

// user
type CreateUserCommand struct {
	UserName string `json:"user_name"`
	Age      int    `json:"age"`
	Amount   int    `json:"amount"`
}
type UpdateUserCommand struct {
	UserName string `json:"user_name"`
	Age      int    `json:"age"`
	Amount   int    `json:"amount"`
	UserId   uint   `json:"user_id"`
}
type DeleteUserCommand struct {
	UserId uint `json:"user_id"`
}
