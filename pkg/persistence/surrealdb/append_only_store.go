package surrealdb

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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
	lastVersion, err := a.numberOfEventsInStream(event.ID.StreamName)
	lastStreamVersion := persistence.StreamVersion(lastVersion)
	if err != nil {
		return fmt.Errorf("error getting number of events in stream: %w", err)
	}
	if lastStreamVersion != event.ID.StreamVersion {
		return &persistence.ErrUnexpectedVersion{Found: lastStreamVersion, Expected: event.ID.StreamVersion}
	}

	id := fmt.Sprintf(`event:{stream_name: '%s', stream_version: %d}`, event.ID.StreamName, event.ID.StreamVersion)
	query := fmt.Sprintf(`
CREATE %s 
SET 
	event_name = $event_name, 
	event_data = $event_data, 
	happened_on = $happened_on;
`, id)
	result, err := a.db.Query(query, map[string]any{
		"event_name":  event.EventName,
		"event_data":  base64.StdEncoding.EncodeToString(event.EventData),
		"happened_on": event.HappenedOn.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("error appending event: %w", err)
	} else if result != nil {
		if strings.Contains(gabs.Wrap(result).Index(0).Path("detail").String(), "already exists") {
			return fmt.Errorf("event already exists")
		}
	}
	return nil
}

func (a *AppendOnlyStore) ReadAllRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `
select 
    id.stream_version as stream_version, 
    id.stream_name as stream_name, 
    event_name, 
    event_data, 
    happened_on 
from event;`

	result, err := a.db.Query(query, nil)
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	return resultToStoredStreamEvent(result)
}

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamName persistence.StreamName) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `
select 
    id.stream_version as stream_version, 
    id.stream_name as stream_name, 
    event_name, 
    event_data, 
    happened_on 
from event 
where 
    id.stream_name = $stream_name;`

	result, err := a.db.Query(query, map[string]any{"stream_name": streamName})
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	return resultToStoredStreamEvent(result)
}

func (a *AppendOnlyStore) ReadEventsByName(ctx context.Context, eventName string) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `
select 
    id.stream_version as stream_version, 
    id.stream_name as stream_name, 
    event_name, 
    event_data, 
    happened_on 
from event 
where 
    event_name = $event_name;`

	result, err := a.db.Query(query, map[string]any{"event_name": eventName})
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	return resultToStoredStreamEvent(result)
}

func resultToStoredStreamEvent(result any) ([]persistence.StoredStreamEvent, error) {
	resultsInQuery := gabs.Wrap(resultFromQuery(result)).Children()
	storedStreamEvents := make([]persistence.StoredStreamEvent, 0, len(resultsInQuery))

	for _, line := range resultsInQuery {
		happenedOn, err := time.Parse(time.RFC3339, line.Path("happened_on").Data().(string))
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %w", err)
		}
		streamVersion, err := strconv.ParseUint(line.Path("stream_version").String(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing stream version: %w", err)
		}

		eventData, err := base64.StdEncoding.DecodeString(line.Path("event_data").Data().(string))
		if err != nil {
			return nil, fmt.Errorf("error decoding base64 string: %w", err)
		}

		storedStreamEvents = append(storedStreamEvents, persistence.StoredStreamEvent{
			ID: persistence.StreamID{
				StreamName:    persistence.StreamName(line.Path("stream_name").Data().(string)),
				StreamVersion: persistence.StreamVersion(streamVersion),
			},
			EventName:  line.Path("event_name").Data().(string),
			EventData:  eventData,
			HappenedOn: happenedOn,
		})
	}

	return storedStreamEvents, nil
}

func (a *AppendOnlyStore) numberOfEventsInStream(streamName persistence.StreamName) (uint64, error) {
	result, err := a.db.Query(fmt.Sprintf(`select id.stream_name, count() from event where id.stream_name = '%s' group by id.stream_name;`, streamName), nil)
	if err != nil {
		return 0, fmt.Errorf("error getting stream version: %w", err)
	}

	count := gabs.Wrap(resultFromQuery(result)).Index(0).Path("count").String()
	parseUint, _ := strconv.ParseUint(count, 10, 64)
	return parseUint, nil
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

	return resultToStoredStreamEvent(records)
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
