package entity

import (
	"github.com/ali-mahdavi-dev/bunny-go/internal/account/domain/events"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/adapter"
)

type User struct {
	adapter.BaseEntity
	AvatarIdentifier string `json:"avatar_identifier" gorm:"avatar_identifier"`
	UserName         string `json:"user_name" gorm:"user_name"`
	FirstName        string `json:"first_name" gorm:"first_name"`
	LastName         string `json:"last_name" gorm:"last_name"`
	Email            string `json:"email" gorm:"email"`
	Password         string `json:"password" gorm:"password"`
}

func NewUser(
	avatarIdentifier string,
	userName string,
	firstName string,
	lastName string,
	email string,
	password string,
) *User {
	return &User{
		AvatarIdentifier: avatarIdentifier,
		UserName:         userName,
		FirstName:        firstName,
		LastName:         lastName,
		Email:            email,
		Password:         password,
		BaseEntity: adapter.BaseEntity{
			Events: []any{
				&events.RegisterUserEvent{
					UserName: userName,
					Email:    email,
				},
			},
		},
	}
}
