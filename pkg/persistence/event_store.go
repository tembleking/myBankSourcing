package persistence

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

type EventStream struct {
	// Name is commonly the aggregate id, but can be any value as long as it is unique for an event stream
	Name string

	// Version is the version of the last event in the stream
	Version uint64

	// Events is the list of events in the stream
	Events []domain.Event
}

type EventStore struct {
	serializer   serializer.EventSerializer
	deserializer serializer.EventDeserializer

	appendOnlyStore AppendOnlyStore
}

// LoadEventStream loads all events for a given aggregate id
func (e *EventStore) LoadEventStream(ctx context.Context, streamName string) (*EventStream, error) {
	return e.LoadEventStreamSubset(ctx, streamName)
}

// LoadEventStreamSubset loads a subset of events for a given aggregate id
// starting from a given version and up to a maximum count
// If the start version is 0, the stream will start from the beginning
// If the max count is 0, the stream will return all events from the start version
func (e *EventStore) LoadEventStreamSubset(ctx context.Context, streamName string) (*EventStream, error) {
	records, err := e.appendOnlyStore.ReadRecords(ctx, streamName)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	stream := &EventStream{
		Name:    streamName,
		Version: uint64(0),
		Events:  make([]domain.Event, 0, len(records)),
	}

	for _, record := range records {
		event, err := e.deserializer.Deserialize(record.Data)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event: %w", err)
		}
		stream.Events = append(stream.Events, event...)
		stream.Version = record.Version
	}

	return stream, nil
}

// AppendToStream appends a list of events to the event stream for a given aggregate id
// returning an error if the expected version does not match the current version
func (e *EventStore) AppendToStream(ctx context.Context, streamName string, expectedVersion uint64, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	data, err := e.serializer.Serialize(events)
	if err != nil {
		return fmt.Errorf("error serializing events: %w", err)
	}

	err = e.appendOnlyStore.Append(ctx, streamName, data, expectedVersion)
	if err != nil {
		return fmt.Errorf("error appending to stream: %w", err)
	}

	return nil
}

func NewEventStore(serializer serializer.EventSerializer, deserializer serializer.EventDeserializer, appendOnlyStore AppendOnlyStore) *EventStore {
	return &EventStore{
		serializer:      serializer,
		deserializer:    deserializer,
		appendOnlyStore: appendOnlyStore,
	}
}
