package entity

import (
	"shikposh-backend/internal/account/domain/events"
	"shikposh-backend/pkg/framework/adapter"
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
	user := &User{
		AvatarIdentifier: avatarIdentifier,
		UserName:         userName,
		FirstName:        firstName,
		LastName:         lastName,
		Email:            email,
		Password:         password,
	}
	user.AddEvent(&events.RegisterUserEvent{
		UserID:           user.ID,
		AvatarIdentifier: user.AvatarIdentifier,
		UserName:         user.UserName,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Email:            user.Email,
	})
	
	return user
}
