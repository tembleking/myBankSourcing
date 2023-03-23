package inmemory

import (
	"context"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AppendOnlyStore struct {
	eventsByStream map[string][]persistence.StoredStreamEvent
	eventsByName   map[string][]persistence.StoredStreamEvent
	rwMutex        sync.RWMutex
}

func (a *AppendOnlyStore) appendEvent(event persistence.StoredStreamEvent) error {
	eventsByStream, _ := a.eventsByStream[event.StreamID]
	eventsByName, _ := a.eventsByName[event.EventName]

	currentVersion := uint64(len(eventsByStream))
	if currentVersion != event.StreamVersion {
		return &persistence.ErrUnexpectedVersion{Found: currentVersion, Expected: event.StreamVersion}
	}

	a.eventsByStream[event.StreamID] = append(eventsByStream, event)
	a.eventsByName[event.EventName] = append(eventsByName, event)
	return nil
}

func (a *AppendOnlyStore) Append(_ context.Context, events ...persistence.StoredStreamEvent) error {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	for _, event := range events {
		if err := a.appendEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamID string) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.eventsByStream[streamID]
	if !ok {
		return nil, &persistence.ErrRecordsNotFoundForStream{StreamID: streamID}
	}

	return fields, nil
}

func (a *AppendOnlyStore) ReadEventsByName(ctx context.Context, eventName string) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.eventsByName[eventName]
	if !ok {
		return nil, &persistence.ErrRecordsNotFoundForEvent{EventName: eventName}
	}

	return fields, nil
}

func NewAppendOnlyStore() *AppendOnlyStore {
	return &AppendOnlyStore{
		eventsByStream: make(map[string][]persistence.StoredStreamEvent),
		eventsByName:   make(map[string][]persistence.StoredStreamEvent),
	}
}
