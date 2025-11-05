package commandeventhandler

import (
	"context"
)

type CommandHandler interface {
	NewCommand() any

	Handle(ctx context.Context, cmd any) (any, error)
}

func NewCommandHandler[Command any](
	handleFunc func(ctx context.Context, cmd *Command) error,
) CommandHandler {
	return &genericCommandHandler[Command, any]{
		handleFunc: func(ctx context.Context, cmd *Command) (any, error) {
			err := handleFunc(ctx, cmd)
			return nil, err
		},
	}
}

func NewCommandHandlerWithResult[Command any, Result any](
	handleFunc func(ctx context.Context, cmd *Command) (Result, error),
) CommandHandler {
	return &genericCommandHandler[Command, Result]{
		handleFunc: func(ctx context.Context, cmd *Command) (any, error) {
			result, err := handleFunc(ctx, cmd)
			return result, err
		},
	}
}

type genericCommandHandler[Command any, Result any] struct {
	handleFunc func(ctx context.Context, cmd *Command) (any, error)
}

func (c genericCommandHandler[Command, Result]) NewCommand() any {
	tVar := new(Command)
	return tVar
}

func (c genericCommandHandler[Command, Result]) Handle(ctx context.Context, cmd any) (any, error) {
	command := cmd.(*Command)
	return c.handleFunc(ctx, command)
}
