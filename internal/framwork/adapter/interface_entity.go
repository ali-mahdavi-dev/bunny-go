package adapter

import (
	"time"

	"gorm.io/gorm"

	commandeventhandler "bunny-go/internal/framwork/service_layer/command_event_handler"
)

type Entity interface {
	GetID() uint
	Event() []commandeventhandler.EventHandler
}

type BaseEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *BaseEntity) GetID() uint {
	return u.ID
}
