package surrealdb

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Jeffail/gabs/v2"
	surreal "github.com/surrealdb/surrealdb.go"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AppendOnlyStore struct {
	db      *surreal.DB
	rwMutex sync.RWMutex
}

func (a *AppendOnlyStore) Append(ctx context.Context, events ...persistence.StoredStreamEvent) error {
	err := a.executeInTransaction(func() error {
		for _, event := range events {
			err := a.appendEvent(event)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (a *AppendOnlyStore) appendEvent(event persistence.StoredStreamEvent) error {
	result, err := a.db.Create("event", storedStreamEventToSurreal(event))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return &persistence.ErrUnexpectedVersion{
				StreamName: event.ID.StreamName,
				Expected:   event.ID.StreamVersion,
			}
		}
		return fmt.Errorf("error appending event: %w", err)
	} else if result != nil {
		if strings.Contains(gabs.Wrap(result).Index(0).Path("detail").String(), "already exists") {
			return &persistence.ErrUnexpectedVersion{
				StreamName: event.ID.StreamName,
				Expected:   event.ID.StreamVersion,
			}
		}
	}
	return nil
}

func (a *AppendOnlyStore) ReadAllRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `select * from event;`
	result, err := a.db.Query(query, nil)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	return resultToStoredStreamEventSlice(result)
}

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamName persistence.StreamName) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `select * from event where id.stream_name = $stream_name;`
	result, err := a.db.Query(query, map[string]any{"stream_name": streamName})
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	return resultToStoredStreamEventSlice(result)
}

func (a *AppendOnlyStore) ReadEventsByName(ctx context.Context, eventName string) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `select * from event where event_name = $event_name;`
	result, err := a.db.Query(query, map[string]any{"event_name": eventName})
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	return resultToStoredStreamEventSlice(result)
}

func resultToStoredStreamEventSlice(result interface{}) ([]persistence.StoredStreamEvent, error) {
	var eventsResponse []surrealStoredStreamEventOut
	_, err := surreal.UnmarshalRaw(result, &eventsResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	storedStreamEvents := make([]persistence.StoredStreamEvent, 0, len(eventsResponse))
	for _, eventResponse := range eventsResponse {
		storedStreamEvents = append(storedStreamEvents, eventResponse.ToStoredStreamEvent())
	}
	return storedStreamEvents, nil
}

func (a *AppendOnlyStore) ReadUndispatchedRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	update_query := `
update event
set reserved_until = time::now() + 5s where (reserved_until is none or reserved_until < time::now()) and (event_dispatched is none or event_dispatched = false)
return 
    id.stream_version as stream_version, 
    id.stream_name as stream_name, 
    event_name, 
    event_data, 
    happened_on;`

	records, err := a.db.Query(update_query, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating records: %w", err)
	}

	return resultToStoredStreamEventSlice(records)
}

func (a *AppendOnlyStore) MarkRecordsAsDispatched(ctx context.Context, streamIDs ...persistence.StreamID) error {
	update_query := `
update event
set event_dispatched = true, reserved_until = none where id.stream_name = $stream_name and id.stream_version = $stream_version
return none;`
	err := a.executeInTransaction(func() error {
		for _, eventID := range streamIDs {
			_, err := a.db.Query(update_query, map[string]any{
				"stream_name":    eventID.StreamName,
				"stream_version": eventID.StreamVersion,
			})
			if err != nil {
				return fmt.Errorf("error updating record: %w", err)
			}
		}
		return nil
	})
	return err
}

func resultFromQuery(result any) any {
	return gabs.Wrap(result).Index(0).Path("result").Data()
}

func (a *AppendOnlyStore) executeInTransaction(f func() error) (err error) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	_, err = a.db.Query("BEGIN TRANSACTION;", nil)
	if err != nil {
		err = fmt.Errorf("error starting transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			if _, errRollback := a.db.Query("CANCEL TRANSACTION;", nil); errRollback != nil {
				err = fmt.Errorf("error rolling back transaction: %w, creation error was %w", errRollback, err)
				return
			}
			return
		}

		_, err = a.db.Query("COMMIT TRANSACTION;", nil)
		if err != nil {
			err = fmt.Errorf("error committing transaction: %w", err)
		}
	}()

	err = f()
	return
}

func NewAppendOnlyStore(db *surreal.DB) *AppendOnlyStore {
	return &AppendOnlyStore{
		db: db,
	}
}
