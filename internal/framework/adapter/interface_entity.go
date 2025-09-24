package adapter

import (
	"time"

	"gorm.io/gorm"
)

type Entity interface {
	GetID() uint64
	Event() []any
}

type BaseEntity struct {
	ID        uint64 `gorm:"primaryKey"`
	Events    []any  `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *BaseEntity) GetID() uint64 {
	return u.ID
}
func (u *BaseEntity) Event() []any {
	events := append([]any(nil), u.Events...)
	u.Events = nil
	return events
}
