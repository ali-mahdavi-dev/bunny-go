package adapter

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Entity interface {
	GetID() uint64
	Event() []any
	AddEvent(event any)
}

type BaseEntity struct {
	ID        uint64 `gorm:"primaryKey"`
	Events    []any  `gorm:"-"`
	eventsMu  sync.Mutex
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *BaseEntity) GetID() uint64 {
	return u.ID
}

// Event returns all events and clears them atomically
func (u *BaseEntity) Event() []any {
	u.eventsMu.Lock()
	defer u.eventsMu.Unlock()
	
	events := append([]any(nil), u.Events...)
	u.Events = nil
	return events
}

// AddEvent adds an event to the entity in a thread-safe manner
func (u *BaseEntity) AddEvent(event any) {
	u.eventsMu.Lock()
	defer u.eventsMu.Unlock()
	u.Events = append(u.Events, event)
}
