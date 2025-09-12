package adapter

import (
	"time"

	"gorm.io/gorm"

	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/command_event_handler"
)

type Entity interface {
	GetID() uint
	Event() []commandeventhandler.EventHandler
}

type BaseEntity struct {
	ID        uint                               `gorm:"primaryKey"`
	Events    []commandeventhandler.EventHandler `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *BaseEntity) GetID() uint {
	return u.ID
}
func (u *BaseEntity) Event() []commandeventhandler.EventHandler {
	events := u.Events
	u.Events = []commandeventhandler.EventHandler{}
	return events
}