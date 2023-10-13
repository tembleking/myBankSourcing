package persistence

import (
	"context"
	"sync"
	"time"
)

type trackableStoredStreamEvent struct {
	StoredStreamEvent
	reservedForDispatch bool
}

type InMemoryAppendOnlyStore struct {
	allEvents          []StoredStreamEvent
	eventsByStream     map[StreamName][]StoredStreamEvent
	eventsByName       map[string][]StoredStreamEvent
	undispatchedEvents []*trackableStoredStreamEvent

	rwMutex sync.RWMutex
}

func (a *InMemoryAppendOnlyStore) appendEvent(event StoredStreamEvent) error {
	eventsByStream := a.eventsByStream[event.ID.StreamName]
	eventsByName := a.eventsByName[event.EventName]

	currentVersion := StreamVersion(len(eventsByStream))
	if currentVersion != event.ID.StreamVersion {
		return &ErrUnexpectedVersion{StreamName: event.ID.StreamName, Expected: event.ID.StreamVersion}
	}

	a.allEvents = append(a.allEvents, event)
	a.eventsByStream[event.ID.StreamName] = append(eventsByStream, event)
	a.eventsByName[event.EventName] = append(eventsByName, event)
	return nil
}

func (a *InMemoryAppendOnlyStore) Append(_ context.Context, events ...StoredStreamEvent) error {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	for _, event := range events {
		if err := a.appendEvent(event); err != nil {
			return err
		}
		a.undispatchedEvents = append(a.undispatchedEvents, &trackableStoredStreamEvent{StoredStreamEvent: event})
	}

	return nil
}

func (a *InMemoryAppendOnlyStore) ReadAllRecords(_ context.Context) ([]StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	return a.allEvents, nil
}

func (a *InMemoryAppendOnlyStore) ReadRecords(_ context.Context, streamID StreamName) ([]StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.eventsByStream[streamID]
	if !ok {
		return nil, nil
	}

	return fields, nil
}

func (a *InMemoryAppendOnlyStore) ReadEventsByName(_ context.Context, eventName string) ([]StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.eventsByName[eventName]
	if !ok {
		return nil, nil
	}

	return fields, nil
}

func (a *InMemoryAppendOnlyStore) ReadUndispatchedRecords(_ context.Context) ([]StoredStreamEvent, error) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	undispachedEvents := make([]StoredStreamEvent, 0, len(a.undispatchedEvents))
	for _, event := range a.undispatchedEvents {
		if !event.reservedForDispatch {
			undispachedEvents = append(undispachedEvents, event.StoredStreamEvent)
			event.reservedForDispatch = true
		}
	}
	go a.unreserveAfterSomeTime(a.undispatchedEvents)
	return undispachedEvents, nil
}

func (a *InMemoryAppendOnlyStore) unreserveAfterSomeTime(events []*trackableStoredStreamEvent) {
	<-time.After(5 * time.Second)
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	for _, event := range events {
		event.reservedForDispatch = false
	}
}

func (a *InMemoryAppendOnlyStore) MarkRecordsAsDispatched(ctx context.Context, events ...StreamID) error {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	undispatchedEvents := make([]*trackableStoredStreamEvent, 0, len(a.undispatchedEvents))
	for _, undispatchedEvent := range a.undispatchedEvents {
		if isEventInList(undispatchedEvent.StoredStreamEvent, events) {
			continue
		}
		undispatchedEvents = append(undispatchedEvents, undispatchedEvent)
	}
	a.undispatchedEvents = undispatchedEvents
	return nil
}

func isEventInList(event StoredStreamEvent, events []StreamID) bool {
	for _, e := range events {
		if e.Equal(event.ID) {
			return true
		}
	}
	return false
}

func NewInMemoryAppendOnlyStore() *InMemoryAppendOnlyStore {
	return &InMemoryAppendOnlyStore{
		allEvents:      make([]StoredStreamEvent, 0),
		eventsByStream: make(map[StreamName][]StoredStreamEvent),
		eventsByName:   make(map[string][]StoredStreamEvent),
	}
}
