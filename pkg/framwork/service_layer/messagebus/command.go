package messagebus

import (
	"context"
)

type CommandHandler interface {
	NewCommand() any

	Handle(ctx context.Context, cmd any) error
}

func NewCommandHandler[Command any](
	handleFunc func(ctx context.Context, cmd *Command) error,
) CommandHandler {
	return &genericCommandHandler[Command]{
		handleFunc: handleFunc,	
	}
}

type genericCommandHandler[Command any] struct {
	handleFunc func(ctx context.Context, cmd *Command) error
}

func (c genericCommandHandler[Command]) NewCommand() any {
	tVar := new(Command)
	return tVar
}

func (c genericCommandHandler[Command]) Handle(ctx context.Context, cmd any) error {
	command := cmd.(*Command)
	return c.handleFunc(ctx, command)
}
