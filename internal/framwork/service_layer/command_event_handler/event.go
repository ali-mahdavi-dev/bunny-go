package commandeventhandler

import (
	"context"
)

type EventHandler interface {
	NewEvent() any

	Handle(ctx context.Context, cmd any) error
}

func NewEventHandler[Event any](
	handleFunc func(ctx context.Context, cmd *Event) error,
) EventHandler {
	return &genericEventHandler[Event]{
		handleFunc: handleFunc,
	}
}

type genericEventHandler[Event any] struct {
	handleFunc func(ctx context.Context, cmd *Event) error
}


func (c genericEventHandler[Event]) NewEvent() any {
	tVar := new(Event)
	return tVar
}

func (c genericEventHandler[Event]) Handle(ctx context.Context, event any) error {
	e := event.(*Event)
	return c.handleFunc(ctx, e)
}
