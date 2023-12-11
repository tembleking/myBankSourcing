package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/golang-migrate/migrate/v4"
	migratesqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite/internal/sqlgen/model"
	. "github.com/tembleking/myBankSourcing/pkg/persistence/sqlite/internal/sqlgen/table"
)

const (
	// https://www.sqlite.org/rescode.html#constraint
	sqliteErrorConstraint = 19
	// https://www.sqlite.org/rescode.html#constraint_unique
	sqliteErrorConstraintUnique = 2067
)

type AppendOnlyStore struct {
	db *sql.DB
}

func (a *AppendOnlyStore) Append(ctx context.Context, events ...persistence.StoredStreamEvent) error {
	if len(events) == 0 {
		return nil
	}

	insertStmt := Event.INSERT(Event.StreamName, Event.StreamVersion, Event.EventName, Event.EventData, Event.HappenedOn, Event.ContentType)
	for _, event := range events {
		insertStmt.MODEL(model.Event{
			StreamName:    string(event.ID.StreamName),
			StreamVersion: int32(event.ID.StreamVersion),
			EventName:     event.EventName,
			EventData:     event.EventData,
			HappenedOn:    event.HappenedOn,
			ContentType:   event.ContentType,
		})
	}

	err := a.doAtomically(ctx, func(tx *sql.Tx) error {
		_, err := insertStmt.ExecContext(ctx, tx)
		return err
	})
	if isErrorUniqueConstraintViolation(err) {
		return persistence.ErrUnexpectedVersion
	}
	if err != nil {
		return fmt.Errorf("unable to push stored stream event into the sqlite append only store: %w", err)
	}
	return nil
}

func isErrorUniqueConstraintViolation(err error) bool {
	var sqliteError sqlite3.Error
	return errors.As(err, &sqliteError) &&
		sqliteError.Code == sqliteErrorConstraint &&
		sqliteError.ExtendedCode == sqliteErrorConstraintUnique
}

func (a *AppendOnlyStore) ReadAllRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	return a.readRecodsWithQuery(ctx, Event.SELECT(Event.AllColumns))
}

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamName persistence.StreamName) ([]persistence.StoredStreamEvent, error) {
	return a.readRecodsWithQuery(ctx, Event.SELECT(Event.AllColumns).WHERE(Event.StreamName.EQ(String(string(streamName)))))
}

func (a *AppendOnlyStore) readRecodsWithQuery(ctx context.Context, query SelectStatement) ([]persistence.StoredStreamEvent, error) {
	var dbEvents []model.Event
	err := query.QueryContext(ctx, a.db, &dbEvents)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve records from stream: %w", err)
	}

	var events []persistence.StoredStreamEvent
	for _, event := range dbEvents {
		events = append(events, modelEventToPersistence(event))
	}
	return events, nil
}

func modelEventToPersistence(dbEvent model.Event) persistence.StoredStreamEvent {
	return persistence.StoredStreamEvent{
		ID: persistence.StreamID{
			StreamName:    persistence.StreamName(dbEvent.StreamName),
			StreamVersion: persistence.StreamVersion(dbEvent.StreamVersion),
		},
		EventName:   dbEvent.EventName,
		EventData:   dbEvent.EventData,
		HappenedOn:  dbEvent.HappenedOn,
		ContentType: dbEvent.ContentType,
	}
}

//go:embed internal/migrations
var migrations embed.FS

func (a *AppendOnlyStore) MigrateDB(ctx context.Context) (err error) {
	driver, err := migratesqlite3.WithInstance(a.db, &migratesqlite3.Config{})
	if err != nil {
		return fmt.Errorf("unable to create migration driver: %w", err)
	}

	fs, err := iofs.New(migrations, "internal/migrations")
	if err != nil {
		return fmt.Errorf("unable to create migration fs: %w", err)
	}

	instance, err := migrate.NewWithInstance("iofs", fs, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration instance: %w", err)
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to migrate: %w", err)
	}

	return nil
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

func InMemory() *AppendOnlyStore {
	db, err := New(":memory:")
	if err != nil {
		panic(fmt.Errorf("this should not have happened: %w", err))
	}

	err = db.MigrateDB(context.Background())
	if err != nil {
		panic(fmt.Errorf("this should not have happened: %w", err))
	}

	return db
}

func (a *AppendOnlyStore) doAtomically(ctx context.Context, function func(tx *sql.Tx) error) (err error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		err = fmt.Errorf("unable to begin transaction: %w", err)
		return
	}

	defer func() {
		if panicMsg := recover(); panicMsg != nil {
			_ = tx.Rollback()
			panic(panicMsg)
		}

		if err != nil {
			rollbackErr := tx.Rollback()
			err = errors.Join(err, rollbackErr)
			return
		}
		err = tx.Commit()
	}()

	err = function(tx)
	return
}

func (a *AppendOnlyStore) Close() error {
	return a.db.Close()
}
