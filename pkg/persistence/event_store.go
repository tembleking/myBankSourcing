package persistence

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

type StreamEvent struct {
	// StreamID is commonly the aggregate id, but can be any value as long as it is unique for an event stream
	StreamID string

	// StreamVersion is the version of the last event in the stream
	StreamVersion uint64

	// Event is the deserialized event
	Event domain.Event
}

type EventStore struct {
	serializer   serializer.EventSerializer
	deserializer serializer.EventDeserializer

	appendOnlyStore AppendOnlyStore
}

// LoadEventStream loads all events for a given aggregate id
func (e *EventStore) LoadEventStream(ctx context.Context, streamName string) ([]StreamEvent, error) {
	records, err := e.appendOnlyStore.ReadRecords(ctx, streamName)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	events := make([]StreamEvent, 0, len(records))
	for _, record := range records {
		event, err := e.deserializer.Deserialize(record.EventData)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event: %w", err)
		}
		events = append(events, StreamEvent{
			StreamID:      streamName,
			StreamVersion: record.StreamVersion,
			Event:         event,
		})
	}

	return events, nil
}

// AppendToStream appends a list of events to the event stream for a given aggregate id
// returning an error if the expected version does not match the current version
func (e *EventStore) AppendToStream(ctx context.Context, streamID string, lastExpectedVersionAfterEventsApplied uint64, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	streamEvents := make([]StoredStreamEvent, 0, len(events))
	version := lastExpectedVersionAfterEventsApplied - uint64(len(events))
	for _, event := range events {
		eventData, err := e.serializer.Serialize(event)
		if err != nil {
			return fmt.Errorf("error serializing event: %w", err)
		}

		streamEvents = append(streamEvents, StoredStreamEvent{
			StreamID:      streamID,
			StreamVersion: version,
			EventName:     event.EventName(),
			EventData:     eventData,
		})

		version++
	}

	err := e.appendOnlyStore.Append(ctx, streamEvents...)
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
