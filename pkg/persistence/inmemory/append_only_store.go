package inmemory

import (
	"context"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AppendOnlyStore struct {
	fields  map[string][]persistence.StoredStreamEvent
	rwMutex sync.RWMutex
}

func (a *AppendOnlyStore) appendEvent(event persistence.StoredStreamEvent) error {
	streamEvents, ok := a.fields[event.StreamID]
	if !ok {
		streamEvents = make([]persistence.StoredStreamEvent, 0)
	}

	currentVersion := uint64(len(streamEvents))
	if currentVersion != event.StreamVersion {
		return &persistence.ErrUnexpectedVersion{Found: currentVersion, Expected: event.StreamVersion}
	}

	a.fields[event.StreamID] = append(streamEvents, event)
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

	fields, ok := a.fields[streamID]
	if !ok {
		return nil, &persistence.ErrRecordsNotFound{StreamID: streamID}
	}

	return fields, nil
}

func NewAppendOnlyStore() *AppendOnlyStore {
	return &AppendOnlyStore{
		fields: make(map[string][]persistence.StoredStreamEvent),
	}
}
