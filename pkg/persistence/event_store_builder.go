package persistence

import (
	"github.com/tembleking/myBankSourcing/pkg/clock"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

type EventStoreBuilder struct {
	serializer      DomainEventSerializer
	deserializer    DomainEventDeserializer
	appendOnlyStore AppendOnlyStore
	clock           Clock
}

func NewEventStoreBuilder(appendOnlyStore AppendOnlyStore) *EventStoreBuilder {
	defaultSerializer := &serializer.GoBinarySerializer{}
	defaultDeserializer := &serializer.GoBinarySerializer{}

	defaultClock := clock.System{}

	return &EventStoreBuilder{
		serializer:      defaultSerializer,
		deserializer:    defaultDeserializer,
		appendOnlyStore: appendOnlyStore,
		clock:           defaultClock,
	}
}

func (b *EventStoreBuilder) WithSerializer(serializer DomainEventSerializer) *EventStoreBuilder {
	b.serializer = serializer
	return b
}

func (b *EventStoreBuilder) WithDeserializer(deserializer DomainEventDeserializer) *EventStoreBuilder {
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
	if b.appendOnlyStore == nil {
		panic("append only store type not set")
	}

	return &EventStore{
		serializer:      b.serializer,
		deserializer:    b.deserializer,
		appendOnlyStore: b.appendOnlyStore,
		clock:           b.clock,
	}
}
