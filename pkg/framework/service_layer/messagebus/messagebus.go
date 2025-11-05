package messagebus

import (
	"context"
	"fmt"
	"reflect"

	"shikposh-backend/pkg/framework/infrastructure/logging"
	commandeventhandler "shikposh-backend/pkg/framework/service_layer/command_event_handler"
	"shikposh-backend/pkg/framework/service_layer/unit_of_work"
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
	AddHandlerEvent(handlers ...commandeventhandler.EventHandler) error
	Handle(ctx context.Context, cmd any) (any, error)
}

type messageBus struct {
	handledCommands map[any]commandeventhandler.CommandHandler
	handledEvent    map[any]commandeventhandler.EventHandler
	uow             unit_of_work.PGUnitOfWork
	eventCh         chan any
}

func NewMessageBus(uow unit_of_work.PGUnitOfWork) MessageBus {
	bus := &messageBus{
		handledCommands: make(map[any]commandeventhandler.CommandHandler),
		handledEvent:    make(map[any]commandeventhandler.EventHandler),
		uow:             uow,
		eventCh:         make(chan any, 100),
	}

	// start event handler
	go func(mb *messageBus, evCh chan any) {
		// TODO: Implement proper shutdown mechanism to close channel gracefully
		defer close(evCh)
		for event := range evCh {
			go func(ev any) {
				if err := mb.HandleEvent(context.Background(), ev); err != nil {
					logging.Error("Failed to handle event").WithError(err).Log()
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

func (m *messageBus) AddHandlerEvent(handlers ...commandeventhandler.EventHandler) error {
	for _, handler := range handlers {
		cmdName := reflect.TypeOf(handler.NewEvent()).String()
		if _, ok := m.handledCommands[cmdName]; ok {
			return DuplicateCommandHandlerError{reflect.TypeOf(cmdName).String()}
		}
		m.handledEvent[cmdName] = handler
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

func (m *messageBus) Handle(ctx context.Context, cmd any) (any, error) {
	cmdName := reflect.TypeOf(cmd).String()

	logging.Debug("Handling command").
		WithAny("command_name", cmdName).
		Log()

	if _, ok := m.handledCommands[cmdName]; !ok {
		err := DoesNotExistCommandHandlerError{cmdName}
		logging.Error("Command handler not found").
			WithAny("command_name", cmdName).
			WithError(err).
			Log()
		return nil, err

	}

	result, err := m.handledCommands[cmdName].Handle(ctx, cmd)
	if err != nil {
		logging.Error("Command handler failed").
			WithAny("command_name", cmdName).
			WithError(err).
			Log()
		return nil, err

	}

	logging.Debug("Collecting domain events from transaction").
		WithAny("command_name", cmdName).
		Log()

	// collect new events from the transaction
	m.uow.CollectNewEvents(ctx, m.eventCh)

	logging.Debug("Command handled successfully").
		WithAny("command_name", cmdName).
		Log()

	return result, nil
}

func (m *messageBus) HandleEvent(ctx context.Context, event any) error {
	eventName := reflect.TypeOf(event).String()
	logging.Debug("Handling event").
		WithAny("event_name", eventName).
		Log()

	if _, ok := m.handledEvent[eventName]; !ok {
		err := DoesNotExistEventHandlerError{reflect.TypeOf(eventName).String()}
		logging.Error("Event handler not found").
			WithAny("event_name", eventName).
			WithError(err).
			Log()
		return err
	}

	err := m.handledEvent[eventName].Handle(ctx, event)
	if err != nil {
		logging.Error("Event handler failed").
			WithAny("event_name", eventName).
			WithError(err).
			Log()
		return err
	}

	logging.Debug("Event handled successfully").
		WithAny("event_name", eventName).
		Log()

	return nil
}
