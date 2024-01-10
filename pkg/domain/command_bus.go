package domain

import "context"

type CommandListener interface {
	OnCommand(ctx context.Context, command Command) error
}

type CommandBus interface {
	Publish(ctx context.Context, commands ...Command) error
	Subscribe(ctx context.Context, listener CommandListener) error
}
