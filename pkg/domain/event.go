package domain

import "github.com/google/uuid"

type Event interface {
	ID() string
	IsDomainEvent()
}

type BaseEvent struct {
	id string
}

func (b *BaseEvent) ID() string {
	return b.id
}

func (b *BaseEvent) IsDomainEvent() {}

func NewBaseEvent() BaseEvent {
	return BaseEvent{
		id: uuid.NewString(),
	}
}
