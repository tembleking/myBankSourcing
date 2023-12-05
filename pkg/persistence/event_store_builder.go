package persistence

import (
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

type EventStoreBuilder struct {
	serializer      DomainEventSerializer
	deserializer    DomainEventDeserializer
	appendOnlyStore AppendOnlyStore
}

func NewEventStoreBuilder(appendOnlyStore AppendOnlyStore) *EventStoreBuilder {
	defaultSerializer := &serializer.GoBinarySerializer{}
	defaultDeserializer := &serializer.GoBinarySerializer{}

	return &EventStoreBuilder{
		serializer:      defaultSerializer,
		deserializer:    defaultDeserializer,
		appendOnlyStore: appendOnlyStore,
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

func (b *EventStoreBuilder) Build() *EventStore {
	if b.appendOnlyStore == nil {
		panic("append only store type not set")
	}

	return &EventStore{
		serializer:      b.serializer,
		deserializer:    b.deserializer,
		appendOnlyStore: b.appendOnlyStore,
	}
}
