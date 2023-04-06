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
	_, err := executeInTransaction[any](a, func() (any, error) {
		for _, event := range events {
			err := a.appendEvent(event)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

func (a *AppendOnlyStore) appendEvent(event persistence.StoredStreamEvent) error {
	lastVersion, err := a.numberOfEventsInStream(event.StreamID)
	if err != nil {
		return fmt.Errorf("error getting number of events in stream: %w", err)
	}
	if lastVersion != event.StreamVersion {
		return &persistence.ErrUnexpectedVersion{Found: lastVersion, Expected: event.StreamVersion}
	}

	id := fmt.Sprintf(`event:{stream_id: '%s', stream_version: %d}`, event.StreamID, event.StreamVersion)
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

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamID string) ([]persistence.StoredStreamEvent, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	query := `
select 
    id.stream_version as stream_version, 
    id.stream_id as stream_id, 
    event_name, 
    event_data, 
    happened_on 
from event 
where 
    id.stream_id = $stream_id;`

	result, err := a.db.Query(query, map[string]any{"stream_id": streamID})
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
    id.stream_id as stream_id, 
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
			StreamID:      line.Path("stream_id").Data().(string),
			StreamVersion: streamVersion,
			EventName:     line.Path("event_name").Data().(string),
			EventData:     eventData,
			HappenedOn:    happenedOn,
		})
	}

	return storedStreamEvents, nil
}

func (a *AppendOnlyStore) numberOfEventsInStream(streamID string) (uint64, error) {
	result, err := a.db.Query(fmt.Sprintf(`select id.stream_id, count() from event where id.stream_id = '%s' group by id.stream_id;`, streamID), nil)
	if err != nil {
		return 0, fmt.Errorf("error getting stream version: %w", err)
	}

	count := gabs.Wrap(resultFromQuery(result)).Index(0).Path("count").String()
	parseUint, _ := strconv.ParseUint(count, 10, 64)
	return parseUint, nil
}

func resultFromQuery(result any) any {
	return gabs.Wrap(result).Index(0).Path("result").Data()
}

func executeInTransaction[T any](a *AppendOnlyStore, f func() (T, error)) (result T, err error) {
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

	result, err = f()
	return
}

func NewAppendOnlyStore(db *surreal.DB) *AppendOnlyStore {
	return &AppendOnlyStore{
		db: db,
	}
}
