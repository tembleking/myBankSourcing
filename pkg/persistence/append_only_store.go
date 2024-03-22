package persistence

import (
	"context"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type AppendOnlyStore interface {
	// Append appends the marshalled events to the store.
	// It returns an error if the expected version does not match the current version.
	Append(ctx context.Context, events ...StoredStreamEvent) error

	ReadOnlyStore
}

type ReadOnlyStore interface {
	// ReadAllRecords reads all events in the store.
	ReadAllRecords(ctx context.Context) ([]StoredStreamEvent, error)

	// ReadRecords reads events within a single Stream by their names.
	ReadRecords(ctx context.Context, streamName string) ([]StoredStreamEvent, error)

	// AfterEventID returns a ReadOnlyStore that only contains events that happened after the given eventID.
	AfterEventID(eventID domain.EventID) ReadOnlyStore

	// Limit returns a ReadOnlyStore that only contains the first n events.
	Limit(limit int) ReadOnlyStore
}

type StreamID struct {
	StreamName    string
	StreamVersion uint64
}

type StoredStreamEvent struct {
	HappenedOn  time.Time
	EventID     domain.EventID
	EventName   string
	ContentType string
	ID          StreamID
	EventData   []byte
}
