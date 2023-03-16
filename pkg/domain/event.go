package domain

import "github.com/google/uuid"

type Event interface {
	ID() string
	IsDomainEvent()
}

type BaseEvent struct {
	EventID string
}

func (b *BaseEvent) ID() string {
	return b.EventID
}

func (b *BaseEvent) IsDomainEvent() {}

func NewBaseEvent() BaseEvent {
	return BaseEvent{
		EventID: uuid.NewString(),
	}
}
