package entity

import (
	"shikposh-backend/pkg/framework/adapter"
)

type Profile struct {
	adapter.BaseEntity
	UserID  uint64 `json:"user_id" gorm:"user_id;uniqueIndex"`
	Bio     string `json:"bio" gorm:"bio"`
	Phone   string `json:"phone" gorm:"phone"`
	Address string `json:"address" gorm:"address"`
}

func NewProfile(userID uint64) *Profile {
	return &Profile{
		UserID: userID,
	}
}

