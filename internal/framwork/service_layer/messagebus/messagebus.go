package messagebus

import (
	"context"
	"fmt"
	"reflect"

	commandeventhandler "github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/command_event_handler"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/unit_of_work"
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

// DoesNotExistCommandHandlerError occurs when a handler with the same name does not exists.
type DoesNotExistEventHandlerError struct {
	EventName string
}

func (d DoesNotExistEventHandlerError) Error() string {
	return fmt.Sprintf("%s was not an Event", d.EventName)
}

type MessageBus interface {
	AddHandler(handlers ...commandeventhandler.CommandHandler) error
	Handle(ctx context.Context, cmd any) error
}

type messageBus struct {
	handledCommands map[any]commandeventhandler.CommandHandler
	handledEvent    map[any]commandeventhandler.EventHandler
	uow             unit_of_work.PGUnitOfWork
	eventCh         chan commandeventhandler.EventHandler
}

func NewMessageBus(uow unit_of_work.PGUnitOfWork, eventCh chan commandeventhandler.EventHandler) MessageBus {
	bus := &messageBus{
		handledCommands: make(map[any]commandeventhandler.CommandHandler),
		handledEvent:    make(map[any]commandeventhandler.EventHandler),
		uow:             uow,
		eventCh:         eventCh,
	}

	// start event handler
	go func(mb *messageBus, evCh <-chan commandeventhandler.EventHandler) {
		for event := range evCh {
			go func(ev any) {
				if err := mb.HandleEvent(context.Background(), ev); err != nil {
					fmt.Println("error handling event:", err)
				}
			}(event)
		}
	}(bus, bus.eventCh)

	return bus
}

func (m *messageBus) AddHandler(handlers ...commandeventhandler.CommandHandler) error {
	for _, handler := range handlers {
		cmdName := reflect.TypeOf(handler.NewCommand()).String()
		if _, ok := m.handledCommands[cmdName]; ok {
			return DuplicateCommandHandlerError{reflect.TypeOf(cmdName).String()}
		}
		m.handledCommands[cmdName] = handler
	}

	return nil
}

func (m *messageBus) AddEvent(handlers ...commandeventhandler.EventHandler) error {
	for _, handler := range handlers {
		eventName := reflect.TypeOf(handler.NewEvent()).String()
		if _, ok := m.handledCommands[eventName]; ok {
			return DuplicateCommandHandlerError{reflect.TypeOf(eventName).String()}
		}
		m.handledEvent[eventName] = handler
	}

	return nil
}

func (m *messageBus) Handle(ctx context.Context, cmd any) error {
	cmdName := reflect.TypeOf(cmd).String()
	if _, ok := m.handledCommands[cmdName]; !ok {
		return DoesNotExistCommandHandlerError{cmdName}
	}

	err := m.handledCommands[cmdName].Handle(ctx, cmd)
	if err != nil {
		return err
	}
	// collect new events
	m.uow.CollectNewEvents(m.eventCh)

	return nil
}

func (m *messageBus) HandleEvent(ctx context.Context, event any) error {
	eventName := reflect.TypeOf(event).String()
	if _, ok := m.handledEvent[eventName]; !ok {
		return DoesNotExistEventHandlerError{reflect.TypeOf(eventName).String()}
	}

	return m.handledEvent[eventName].Handle(ctx, event)
}
