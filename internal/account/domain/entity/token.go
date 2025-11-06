package entity

import (
	"shikposh-backend/pkg/framework/adapter"
)

type Token struct {
	adapter.BaseEntity
	Token  string `json:"token" gorm:"token"`
	UserID uint64 `json:"user_id" gorm:"user_id"`
}

func NewToken(token string, userID uint64) *Token {
	return &Token{
		Token:  token,
		UserID: userID,
	}
}
