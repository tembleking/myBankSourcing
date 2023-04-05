package persistence

import (
	"github.com/tembleking/myBankSourcing/pkg/clock"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

type EventStoreBuilder struct {
	serializer      EventSerializer
	deserializer    EventDeserializer
	appendOnlyStore AppendOnlyStore
	clock           Clock
}

func NewEventStoreBuilder() *EventStoreBuilder {
	defaultSerializer := &serializer.GoBinarySerializer{}
	defaultDeserializer := &serializer.GoBinarySerializer{}
	defaultAppendOnlyStore := NewInMemoryAppendOnlyStore()
	defaultClock := clock.System{}

	return &EventStoreBuilder{
		serializer:      defaultSerializer,
		deserializer:    defaultDeserializer,
		appendOnlyStore: defaultAppendOnlyStore,
		clock:           defaultClock,
	}
}

func (b *EventStoreBuilder) WithSerializer(serializer EventSerializer) *EventStoreBuilder {
	b.serializer = serializer
	return b
}

func (b *EventStoreBuilder) WithDeserializer(deserializer EventDeserializer) *EventStoreBuilder {
	b.deserializer = deserializer
	return b
}

func (b *EventStoreBuilder) WithAppendOnlyStore(appendOnlyStore AppendOnlyStore) *EventStoreBuilder {
	b.appendOnlyStore = appendOnlyStore
	return b
}

func (b *EventStoreBuilder) WithClock(clock Clock) *EventStoreBuilder {
	b.clock = clock
	return b
}

func (b *EventStoreBuilder) Build() *EventStore {
	return &EventStore{
		serializer:      b.serializer,
		deserializer:    b.deserializer,
		appendOnlyStore: b.appendOnlyStore,
		clock:           b.clock,
		dispatchers:     []EventDispatcher{},
	}
}
