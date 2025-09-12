package mocks

import (
	"context"
	"errors"
	"reflect"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/helpers/is"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/service_layer/types"
)

var (
	handlerNotFountError = errors.New("no handler found for command")
	handlerInvalidError  = errors.New("invalid handler type")
)

type FakeMessageBus struct {
	Uow      *FakeUnitOfWork
	handlers map[string]types.HandlerType
}

func NewFakeMessageBus(uow *FakeUnitOfWork) *FakeMessageBus {
	return &FakeMessageBus{
		handlers: make(map[string]types.HandlerType),
		Uow:      uow,
	}
}

func (m *FakeMessageBus) Register(cmd types.Command, handler types.HandlerType) {
	m.handlers[reflect.TypeOf(cmd).String()] = handler
}

func (m *FakeMessageBus) Handle(ctx context.Context, cmd types.Command) (any, error) {
	typeCmd := reflect.TypeOf(cmd)
	if is.Ptr(cmd) {
		typeCmd = typeCmd.Elem()
	}
	typeName := typeCmd.String()
	handler, exists := m.handlers[typeName]
	if !exists {
		return nil, handlerNotFountError
	} else if h, ok := handler.(types.HandlerType); ok {
		return h.Handle(ctx, cmd)
	}

	return nil, handlerInvalidError
}
