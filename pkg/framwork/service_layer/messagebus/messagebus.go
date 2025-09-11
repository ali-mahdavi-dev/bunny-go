package messagebus

import (
	"context"
	"fmt"
	"reflect"
)

// DuplicateCommandHandlerError occurs when a handler with the same name already exists.
type DuplicateCommandHandlerError struct {
	CommandName string
}

func (d DuplicateCommandHandlerError) Error() string {
	return fmt.Sprintf("command handler for command %s already exists", d.CommandName)
}

// DoesNotExistCommandHandlerError occurs when a handler with the same name already exists.
type DoesNotExistCommandHandlerError struct {
	CommandName string
}

func (d DoesNotExistCommandHandlerError) Error() string {
	return fmt.Sprintf("%s was not an Command", d.CommandName)
}

type MessageBus interface {
	AddHandler(handlers ...CommandHandler) error
	Handle(ctx context.Context, cmd any) error
}

type messageBus struct {
	handledCommands map[any]CommandHandler
}

func NewMessageBus() MessageBus {
	return &messageBus{
		handledCommands: make(map[any]CommandHandler),
	}
}

func (m *messageBus) AddHandler(handlers ...CommandHandler) error {
	for _, handler := range handlers {
		cmd := handler.NewCommand()
		if _, ok := m.handledCommands[cmd]; ok {
			return DuplicateCommandHandlerError{reflect.TypeOf(cmd).String()}
		}
		m.handledCommands[cmd] = handler
	}

	return nil
}

func (m *messageBus) Handle(ctx context.Context, cmd any) error {
	if _, ok := m.handledCommands[cmd]; !ok {
		return DoesNotExistCommandHandlerError{reflect.TypeOf(cmd).String()}
	}

	return m.handledCommands[cmd].Handle(ctx, cmd)
}
