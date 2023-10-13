package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

const (
	// https://www.sqlite.org/rescode.html#constraint
	sqliteErrorConstraint = 19
	// https://www.sqlite.org/rescode.html#constraint_primarykey
	sqliteErrorConstraintPrimaryKey = 1555
)

type AppendOnlyStore struct {
	db *sql.DB
}

// Append implements persistence.AppendOnlyStore.
func (a *AppendOnlyStore) Append(ctx context.Context, events ...persistence.StoredStreamEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO events (stream_name, stream_version, event_name, event_data, happened_on) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("error creating the prepared statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err = stmt.Exec(event.ID.StreamName, event.ID.StreamVersion, event.EventName, event.EventData, event.HappenedOn)
		var sqliteError sqlite3.Error
		if err != nil {
			if errors.As(err, &sqliteError) && sqliteError.Code == sqliteErrorConstraint && sqliteError.ExtendedCode == sqliteErrorConstraintPrimaryKey {
				return &persistence.ErrUnexpectedVersion{
					StreamName: event.ID.StreamName,
					Expected:   event.ID.StreamVersion,
				}
			}
			return fmt.Errorf("unable to push stored stream event into the sqlite append only store: %w", err)
		}
	}

	return tx.Commit()
}

// MarkRecordsAsDispatched implements persistence.AppendOnlyStore.
func (a *AppendOnlyStore) MarkRecordsAsDispatched(ctx context.Context, events ...persistence.StreamID) error {

	if len(events) == 0 {
		return nil
	}

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE events SET dispatched = true WHERE stream_name = ? AND stream_version = ?;")
	if err != nil {
		return fmt.Errorf("error creating the prepared statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err = stmt.Exec(event.StreamName, event.StreamVersion)
		if err != nil {
			return fmt.Errorf("unable to mark the event as dispatched in the sqlite append only store: %w", err)
		}
	}

	return tx.Commit()
}

// ReadAllRecords implements persistence.AppendOnlyStore.
func (a *AppendOnlyStore) ReadAllRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {

	rows, err := a.db.QueryContext(ctx, `SELECT stream_name, stream_version, event_name, event_data, happened_on FROM events;`)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve records from stream: %w", err)
	}

	events := []persistence.StoredStreamEvent{}
	for rows.Next() {
		event := persistence.StoredStreamEvent{
			ID: persistence.StreamID{},
		}
		err := rows.Scan(
			&event.ID.StreamName,
			&event.ID.StreamVersion,
			&event.EventName,
			&event.EventData,
			&event.HappenedOn,
		)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}
	return events, nil
}

// ReadEventsByName implements persistence.AppendOnlyStore.
func (a *AppendOnlyStore) ReadEventsByName(ctx context.Context, eventName string) ([]persistence.StoredStreamEvent, error) {
	rows, err := a.db.QueryContext(ctx, `SELECT stream_name, stream_version, event_name, event_data, happened_on FROM events WHERE event_name = ?;`, eventName)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve records from stream: %w", err)
	}

	events := []persistence.StoredStreamEvent{}
	for rows.Next() {
		event := persistence.StoredStreamEvent{
			ID: persistence.StreamID{},
		}
		err := rows.Scan(
			&event.ID.StreamName,
			&event.ID.StreamVersion,
			&event.EventName,
			&event.EventData,
			&event.HappenedOn,
		)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}
	return events, nil
}

// ReadRecords implements persistence.AppendOnlyStore.
func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamName persistence.StreamName) ([]persistence.StoredStreamEvent, error) {
	rows, err := a.db.QueryContext(ctx, `SELECT stream_name, stream_version, event_name, event_data, happened_on FROM events WHERE stream_name = ?;`, streamName)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve records from stream: %w", err)
	}

	events := []persistence.StoredStreamEvent{}
	for rows.Next() {
		event := persistence.StoredStreamEvent{
			ID: persistence.StreamID{},
		}
		err := rows.Scan(
			&event.ID.StreamName,
			&event.ID.StreamVersion,
			&event.EventName,
			&event.EventData,
			&event.HappenedOn,
		)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}
	return events, nil
}

// ReadUndispatchedRecords implements persistence.AppendOnlyStore.
func (a *AppendOnlyStore) ReadUndispatchedRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	rows, err := a.db.QueryContext(ctx, `
UPDATE events SET reserved_until = DATETIME(CURRENT_TIMESTAMP, '+5 seconds')
WHERE
	reserved_until <= CURRENT_TIMESTAMP AND dispatched = FALSE
RETURNING stream_name, stream_version, event_name, event_data, happened_on
`)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve records from stream: %w", err)
	}

	events := []persistence.StoredStreamEvent{}
	for rows.Next() {
		event := persistence.StoredStreamEvent{
			ID: persistence.StreamID{},
		}
		err := rows.Scan(
			&event.ID.StreamName,
			&event.ID.StreamVersion,
			&event.EventName,
			&event.EventData,
			&event.HappenedOn,
		)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}
	return events, nil
}

func New(connectionString string) (*AppendOnlyStore, error) {
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open sqlite database connection: %w", err)
	}

	return &AppendOnlyStore{
		db: db,
	}, nil
}

func (s *AppendOnlyStore) MigrateDB(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS events (
	stream_name TEXT NOT NULL, 
	stream_version INTEGER NOT NULL, 
	event_name TEXT NOT NULL, 
	event_data BLOB NOT NULL, 
	happened_on	TIMESTAMP NOT NULL,
	dispatched INTEGER DEFAULT FALSE,
	reserved_until TEXT DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT events_stream_pk PRIMARY KEY (stream_name, stream_version)
);
`)
	if err != nil {
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	return tx.Commit()
}

func (s *AppendOnlyStore) Close() error {
	return s.db.Close()
}
