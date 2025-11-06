package messagebus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	apperrors "shikposh-backend/pkg/framework/errors"
	"shikposh-backend/pkg/framework/infrastructure/logging"
	commandeventhandler "shikposh-backend/pkg/framework/service_layer/command_event_handler"
	"shikposh-backend/pkg/framework/service_layer/unit_of_work"
)

type MessageBus interface {
	AddHandler(handlers ...commandeventhandler.CommandHandler) error
	AddHandlerEvent(handlers ...commandeventhandler.EventHandler) error
	Handle(ctx context.Context, cmd any) (any, error)
	Shutdown(ctx context.Context) error
}

type messageBus struct {
	handledCommands map[any]commandeventhandler.CommandHandler
	handledEvent    map[any]commandeventhandler.EventHandler
	uow             unit_of_work.PGUnitOfWork
	eventCh         chan any
	shutdownCh      chan struct{}
	wg              sync.WaitGroup
	mu              sync.RWMutex
	closed          bool
}

func NewMessageBus(uow unit_of_work.PGUnitOfWork) MessageBus {
	bus := &messageBus{
		handledCommands: make(map[any]commandeventhandler.CommandHandler),
		handledEvent:    make(map[any]commandeventhandler.EventHandler),
		uow:             uow,
		eventCh:         make(chan any, 100),
		shutdownCh:      make(chan struct{}),
		closed:          false,
	}

	// start event handler worker pool
	// Using a single goroutine to process events sequentially to avoid race conditions
	// Events are processed one at a time to ensure proper ordering and prevent concurrency issues
	bus.wg.Add(1)
	go func(mb *messageBus, evCh chan any) {
		defer mb.wg.Done()
		defer close(evCh)
		
		for {
			select {
			case event, ok := <-evCh:
				if !ok {
					logging.Debug("Event channel closed, stopping event handler").Log()
					return
				}
				// Process events sequentially to avoid race conditions
				// Each event handler may create new events, which will be collected and processed
				if err := mb.HandleEvent(context.Background(), event); err != nil {
					logging.Error("Failed to handle event").WithError(err).Log()
				}
			case <-mb.shutdownCh:
				logging.Debug("Shutdown signal received, processing remaining events").Log()
				// Process remaining events before shutdown
				for {
					select {
					case event, ok := <-evCh:
						if !ok {
							return
						}
						if err := mb.HandleEvent(context.Background(), event); err != nil {
							logging.Error("Failed to handle event during shutdown").WithError(err).Log()
						}
					default:
						return
					}
				}
			}
		}
	}(bus, bus.eventCh)

	return bus
}

func (m *messageBus) AddHandler(handlers ...commandeventhandler.CommandHandler) error {
	for _, handler := range handlers {
		cmdName := reflect.TypeOf(handler.NewCommand()).String()
		if _, ok := m.handledCommands[cmdName]; ok {
			return apperrors.Conflict("", fmt.Sprintf("command handler for %s already exists", cmdName))
		}
		m.handledCommands[cmdName] = handler
	}

	return nil
}

func (m *messageBus) AddHandlerEvent(handlers ...commandeventhandler.EventHandler) error {
	for _, handler := range handlers {
		eventName := reflect.TypeOf(handler.NewEvent()).String()
		if _, ok := m.handledEvent[eventName]; ok {
			return apperrors.Conflict("", fmt.Sprintf("event handler for %s already exists", eventName))
		}
		m.handledEvent[eventName] = handler
	}

	return nil
}

func (m *messageBus) AddEvent(handlers ...commandeventhandler.EventHandler) error {
	for _, handler := range handlers {
		eventName := reflect.TypeOf(handler.NewEvent()).String()
		if _, ok := m.handledEvent[eventName]; ok {
			return apperrors.Conflict("", fmt.Sprintf("event handler for %s already exists", eventName))
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

	handler, ok := m.handledCommands[cmdName]
	if !ok {
		err := apperrors.NotFound("", fmt.Sprintf("command handler for %s not found", cmdName))
		logging.Error("Command handler not found").
			WithAny("command_name", cmdName).
			WithError(err).
			Log()
		return nil, err
	}

	result, err := handler.Handle(ctx, cmd)
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
	m.mu.RLock()
	closed := m.closed
	m.mu.RUnlock()
	
	if !closed {
		m.uow.CollectNewEvents(ctx, m.eventCh)
	}

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
		err := apperrors.NotFound("", fmt.Sprintf("event handler for %s not found", eventName))
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

	// Collect and handle nested events that may have been created by the event handler
	// This supports event handlers that create new events
	logging.Debug("Collecting nested events from event handler").
		WithAny("event_name", eventName).
		Log()
	
	m.mu.RLock()
	closed := m.closed
	m.mu.RUnlock()
	
	if !closed {
		m.uow.CollectNewEvents(ctx, m.eventCh)
	}

	logging.Debug("Event handled successfully").
		WithAny("event_name", eventName).
		Log()

	return nil
}

// Shutdown gracefully shuts down the message bus.
// It stops accepting new events, processes remaining events, and closes the event channel.
func (m *messageBus) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	m.mu.Unlock()

	logging.Debug("Initiating message bus shutdown").Log()

	// Signal shutdown to the event handler goroutine
	close(m.shutdownCh)

	// Wait for all events to be processed or context timeout
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logging.Debug("Message bus shutdown completed successfully").Log()
		return nil
	case <-ctx.Done():
		logging.Warn("Message bus shutdown timed out").WithError(ctx.Err()).Log()
		return ctx.Err()
	}
}
