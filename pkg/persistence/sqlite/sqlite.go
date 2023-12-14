package sqlite

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite/internal/model"
)

const (
	// https://www.sqlite.org/rescode.html#constraint
	sqliteErrorConstraint = 19
	// https://www.sqlite.org/rescode.html#constraint_unique
	sqliteErrorConstraintUnique = 2067
)

type AppendOnlyStore struct {
	db *gorm.DB
}

func (a *AppendOnlyStore) Append(ctx context.Context, events ...persistence.StoredStreamEvent) error {
	if len(events) == 0 {
		return nil
	}

	eventsToInsert := make([]model.Event, 0, len(events))
	for _, event := range events {
		eventsToInsert = append(eventsToInsert, model.Event{
			StreamName:    event.ID.StreamName,
			StreamVersion: strconv.FormatUint(event.ID.StreamVersion, 10),
			EventName:     event.EventName,
			EventID:       event.EventID,
			EventData:     event.EventData,
			HappenedOn:    event.HappenedOn,
			ContentType:   event.ContentType,
		})
	}

	err := a.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Omit("row_id").CreateInBatches(eventsToInsert, 1000).Error
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
	return readRecodsWithQuery(ctx, a.db.WithContext(ctx))
}

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, streamName string) ([]persistence.StoredStreamEvent, error) {
	return readRecodsWithQuery(ctx, a.db.WithContext(ctx).Where("stream_name = ?", streamName))
}

func readRecodsWithQuery(ctx context.Context, db *gorm.DB) ([]persistence.StoredStreamEvent, error) {
	var dbEvents []model.Event
	err := db.Find(&dbEvents).Error
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve records from stream: %w", err)
	}

	var events []persistence.StoredStreamEvent
	for _, event := range dbEvents {
		storedStreamEvent, err := modelEventToPersistence(event)
		if err != nil {
			return nil, fmt.Errorf("unable to convert model event to persistence event: %w", err)
		}
		events = append(events, storedStreamEvent)
	}
	return events, nil
}

func modelEventToPersistence(dbEvent model.Event) (persistence.StoredStreamEvent, error) {
	streamVersion, err := strconv.ParseUint(dbEvent.StreamVersion, 10, 64)
	if err != nil {
		return persistence.StoredStreamEvent{}, fmt.Errorf("error parsing stream version '%s' to uint64: %w", dbEvent.StreamVersion, err)
	}

	return persistence.StoredStreamEvent{
		ID: persistence.StreamID{
			StreamName:    dbEvent.StreamName,
			StreamVersion: streamVersion,
		},
		EventName:   dbEvent.EventName,
		EventID:     dbEvent.EventID,
		EventData:   dbEvent.EventData,
		HappenedOn:  dbEvent.HappenedOn,
		ContentType: dbEvent.ContentType,
	}, nil
}

func New(connectionString string) (*AppendOnlyStore, error) {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
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

func (a *AppendOnlyStore) Close() error {
	db, err := a.db.DB()
	if err != nil {
		return fmt.Errorf("unable to retrieve database connection: %w", err)
	}
	return db.Close()
}
