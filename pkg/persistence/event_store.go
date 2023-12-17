package persistence

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

// EventStore is a store for events that can be used to load and save domain events.
// It is a wrapper around an AppendOnlyStore that handles serialization and deserialization of events.
// It also handles dispatching events to any registered EventDispatchers.
// It can be constructed using the EventStoreBuilder.
type EventStore struct {
	serializer      DomainEventSerializer
	appendOnlyStore AppendOnlyStore

	*ReadOnlyEventStore
}

type ReadOnlyEventStore struct {
	deserializer  DomainEventDeserializer
	readOnlyStore ReadOnlyStore
}

// LoadEventStream loads all events for a given aggregate id
func (e *ReadOnlyEventStore) LoadEventStream(ctx context.Context, streamName string) ([]domain.Event, error) {
	records, err := e.readOnlyStore.ReadRecords(ctx, streamName)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	events := make([]domain.Event, 0, len(records))
	for _, record := range records {
		event, err := e.deserializer.DeserializeDomainEvent(record.EventName, record.EventData)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event: %w", err)
		}
		events = append(events, event)
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
	events := aggregate.UncommittedEvents()
	if len(events) == 0 {
		return nil, nil
	}

	storedStreamEvents := make([]StoredStreamEvent, 0, len(aggregate.UncommittedEvents()))
	version := aggregate.Version() - uint64(len(events))
	for _, event := range events {
		eventData, err := e.serializer.SerializeDomainEvent(event)
		if err != nil {
			return nil, fmt.Errorf("error serializing event: %w", err)
		}

		storedStreamEvents = append(storedStreamEvents, StoredStreamEvent{
			ID:          StreamID{StreamName: aggregate.ID(), StreamVersion: version},
			EventID:     event.EventID(),
			EventName:   event.EventName(),
			EventData:   eventData,
			HappenedOn:  event.HappenedOn(),
			ContentType: e.serializer.ContentType(),
		})

		version++
	}

	return storedStreamEvents, nil
}

func (e *ReadOnlyEventStore) LoadAllEvents(ctx context.Context) ([]domain.Event, error) {
	records, err := e.readOnlyStore.ReadAllRecords(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	events := make([]domain.Event, 0, len(records))
	for _, record := range records {
		event, err := e.deserializer.DeserializeDomainEvent(record.EventName, record.EventData)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event '%s' for stream '%s' in version '%d': %w", record.EventName, record.ID.StreamName, record.ID.StreamVersion, err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (e *ReadOnlyEventStore) AfterEventID(eventID domain.EventID) *ReadOnlyEventStore {
	return &ReadOnlyEventStore{
		deserializer:  e.deserializer,
		readOnlyStore: e.readOnlyStore.AfterEventID(eventID),
	}
}

func (e *ReadOnlyEventStore) Limit(limit int) *ReadOnlyEventStore {
	return &ReadOnlyEventStore{
		deserializer:  e.deserializer,
		readOnlyStore: e.readOnlyStore.Limit(limit),
	}
}
