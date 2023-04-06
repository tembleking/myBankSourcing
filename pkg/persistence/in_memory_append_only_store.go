package persistence

import (
	"context"
	"sync"
)

type InMemoryAppendOnlyStore struct {
	allEvents      []StoredStreamEvent
	eventsByStream map[string][]StoredStreamEvent
	eventsByName   map[string][]StoredStreamEvent
	rwMutex        sync.RWMutex
}

func (a *InMemoryAppendOnlyStore) appendEvent(event StoredStreamEvent) error {
	eventsByStream := a.eventsByStream[event.StreamID]
	eventsByName := a.eventsByName[event.EventName]

	currentVersion := uint64(len(eventsByStream))
	if currentVersion != event.StreamVersion {
		return &ErrUnexpectedVersion{Found: currentVersion, Expected: event.StreamVersion}
	}

	a.allEvents = append(a.allEvents, event)
	a.eventsByStream[event.StreamID] = append(eventsByStream, event)
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
	}

	return nil
}

func (a *InMemoryAppendOnlyStore) ReadAllRecords(ctx context.Context) ([]StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	return a.allEvents, nil
}

func (a *InMemoryAppendOnlyStore) ReadRecords(ctx context.Context, streamID string) ([]StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.eventsByStream[streamID]
	if !ok {
		return nil, nil
	}

	return fields, nil
}

func (a *InMemoryAppendOnlyStore) ReadEventsByName(ctx context.Context, eventName string) ([]StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.eventsByName[eventName]
	if !ok {
		return nil, nil
	}

	return fields, nil
}

func NewInMemoryAppendOnlyStore() *InMemoryAppendOnlyStore {
	return &InMemoryAppendOnlyStore{
		allEvents:      make([]StoredStreamEvent, 0),
		eventsByStream: make(map[string][]StoredStreamEvent),
		eventsByName:   make(map[string][]StoredStreamEvent),
	}
}
