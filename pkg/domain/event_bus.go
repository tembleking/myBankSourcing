package domain

import "context"

type EventListener interface {
	OnEvent(ctx context.Context, event Event) error
}

type EventBus interface {
	Publish(ctx context.Context, events ...Event) error
	Subscribe(ctx context.Context, listener EventListener) error
}
