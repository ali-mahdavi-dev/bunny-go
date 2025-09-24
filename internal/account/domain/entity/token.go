package entity

import (
	"time"
)

type Token struct {
	ID        uint64 `gorm:"primaryKey"`
	Token     string `json:"token" gorm:"token"`
	UserID    uint64 `json:"user_id" gorm:"user_id"`
	CreatedAt time.Time
}

func NewToken(token string, userID uint64) *Token {
	return &Token{
		Token:  token,
		UserID: userID,
	}
}

func (u *Token) GetID() uint64 {
	return u.ID
}
func (u *Token) Event() []any {
	return nil
}
