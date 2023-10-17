package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Clock interface {
	Now() time.Time
}

type (
	StreamName    string
	StreamVersion uint64
	StreamID      struct {
		// StreamName is commonly the aggregate id, but can be any value as long as it is unique for an event stream
		StreamName StreamName
		// StreamVersion is the version of the last event in the stream
		StreamVersion StreamVersion
	}
)

func (s *StreamID) Equal(other StreamID) bool {
	return s.StreamName == other.StreamName && s.StreamVersion == other.StreamVersion
}

type StreamEvent struct {
	// ID is the id of the event stream
	ID StreamID

	// Event is the deserialized event
	Event domain.Event

	// HappenedOn is the time the event happened
	HappenedOn time.Time
}

// EventStore is a store for events that can be used to load and save domain events.
// It is a wrapper around an AppendOnlyStore that handles serialization and deserialization of events.
// It also handles dispatching events to any registered EventDispatchers.
// It can be constructed using the EventStoreBuilder.
type EventStore struct {
	serializer      DomainEventSerializer
	deserializer    DomainEventDeserializer
	appendOnlyStore AppendOnlyStore
	clock           Clock
}

// LoadEventStream loads all events for a given aggregate id
func (e *EventStore) LoadEventStream(ctx context.Context, streamName StreamName) ([]StreamEvent, error) {
	records, err := e.appendOnlyStore.ReadRecords(ctx, streamName)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	events := make([]StreamEvent, 0, len(records))
	for _, record := range records {
		event, err := e.deserializer.DeserializeDomainEvent(record.EventName, record.EventData)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event: %w", err)
		}
		events = append(events, StreamEvent{
			ID:         record.ID,
			Event:      event,
			HappenedOn: record.HappenedOn,
		})
	}

	return events, nil
}

// AppendToStream appends a list of events to the event stream for a given aggregate id
// returning an error if the expected version does not match the current version
// FIXME: This should only save one aggregate
func (e *EventStore) AppendToStream(ctx context.Context, aggregates ...domain.Aggregate) error {
	storedStreamEvents := []StoredStreamEvent{}

	for _, aggregate := range aggregates {
		eventsFromAggregate, err := e.streamEventsFromAggregate(aggregate)
		if err != nil {
			return fmt.Errorf("error extracting stream events from aggregate: %w", err)
		}

		storedStreamEvents = append(storedStreamEvents, eventsFromAggregate...)
	}

	err := e.appendOnlyStore.Append(ctx, storedStreamEvents...)
	if err != nil {
		return fmt.Errorf("error appending to stream: %w", err)
	}

	return nil
}

func (e *EventStore) streamEventsFromAggregate(aggregate domain.Aggregate) ([]StoredStreamEvent, error) {
	events := aggregate.Events()
	if len(events) == 0 {
		return nil, nil
	}

	storedStreamEvents := make([]StoredStreamEvent, 0, len(aggregate.Events()))
	version := aggregate.Version() - uint64(len(events))
	for _, event := range events {
		eventData, err := e.serializer.SerializeDomainEvent(event)
		if err != nil {
			return nil, fmt.Errorf("error serializing event: %w", err)
		}

		now := e.clock.Now().UTC()
		storedStreamEvents = append(storedStreamEvents, StoredStreamEvent{
			ID:         StreamID{StreamName: StreamName(aggregate.ID()), StreamVersion: StreamVersion(version)},
			EventName:  event.EventName(),
			EventData:  eventData,
			HappenedOn: now,
		})

		version++
	}

	return storedStreamEvents, nil
}

func (e *EventStore) LoadAllEvents(ctx context.Context) ([]StreamEvent, error) {
	records, err := e.appendOnlyStore.ReadAllRecords(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	events := make([]StreamEvent, 0, len(records))
	for _, record := range records {
		event, err := e.deserializer.DeserializeDomainEvent(record.EventName, record.EventData)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event '%s' for stream '%s' in version '%d': %w", record.EventName, record.ID.StreamName, record.ID.StreamVersion, err)
		}
		events = append(events, StreamEvent{
			ID:         record.ID,
			Event:      event,
			HappenedOn: record.HappenedOn,
		})
	}

	return events, nil
}
