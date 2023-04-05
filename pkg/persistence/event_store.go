package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type EventDispatcher interface {
	Dispatch(ctx context.Context, events ...StreamEvent)
}

type Clock interface {
	Now() time.Time
}

type StreamEvent struct {
	// StreamID is commonly the aggregate id, but can be any value as long as it is unique for an event stream
	StreamID string

	// StreamVersion is the version of the last event in the stream
	StreamVersion uint64

	// Event is the deserialized event
	Event      domain.Event
	HappenedOn time.Time
}

type EventStore struct {
	serializer   serializer.EventSerializer
	deserializer serializer.EventDeserializer

	appendOnlyStore AppendOnlyStore
	dispatchers     []EventDispatcher
	clock           Clock
}

// LoadEventStream loads all events for a given aggregate id
func (e *EventStore) LoadEventStream(ctx context.Context, streamID string) ([]StreamEvent, error) {
	records, err := e.appendOnlyStore.ReadRecords(ctx, streamID)
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
			StreamID:      streamID,
			StreamVersion: record.StreamVersion,
			Event:         event,
			HappenedOn:    record.HappenedOn,
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

	storedStreamEvents := make([]StoredStreamEvent, 0, len(events))
	streamEvents := make([]StreamEvent, 0, len(events))
	version := lastExpectedVersionAfterEventsApplied - uint64(len(events))
	for _, event := range events {
		eventData, err := e.serializer.Serialize(event)
		if err != nil {
			return fmt.Errorf("error serializing event: %w", err)
		}

		now := e.clock.Now().UTC()
		storedStreamEvents = append(storedStreamEvents, StoredStreamEvent{
			StreamID:      streamID,
			StreamVersion: version,
			EventName:     event.EventName(),
			EventData:     eventData,
			HappenedOn:    now,
		})

		streamEvents = append(streamEvents, StreamEvent{
			StreamID:      streamID,
			StreamVersion: version,
			Event:         event,
			HappenedOn:    now,
		})

		version++
	}

	err := e.appendOnlyStore.Append(ctx, storedStreamEvents...)
	if err != nil {
		return fmt.Errorf("error appending to stream: %w", err)
	}

	for _, dispatcher := range e.dispatchers {
		dispatcher.Dispatch(ctx, streamEvents...)
	}

	return nil
}

func (e *EventStore) AddDispatchers(dispatchers ...EventDispatcher) {
	e.dispatchers = append(e.dispatchers, dispatchers...)
}

func (e *EventStore) LoadEventsByName(ctx context.Context, eventName string) ([]StreamEvent, error) {
	records, err := e.appendOnlyStore.ReadEventsByName(ctx, eventName)
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
			StreamID:      record.StreamID,
			StreamVersion: record.StreamVersion,
			Event:         event,
			HappenedOn:    record.HappenedOn,
		})
	}

	return events, nil
}

func NewEventStore(serializer serializer.EventSerializer, deserializer serializer.EventDeserializer, appendOnlyStore AppendOnlyStore, clock Clock) *EventStore {
	return &EventStore{
		serializer:      serializer,
		deserializer:    deserializer,
		appendOnlyStore: appendOnlyStore,
		clock:           clock,
	}
}
