package persistence

import (
	"context"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type AppendOnlyStore interface {
	// Append appends the marshalled events to the store.
	// It returns an error if the expected version does not match the current version.
	Append(ctx context.Context, events ...StoredStreamEvent) error

	// ReadAllRecords reads all events in the store.
	ReadAllRecords(ctx context.Context) ([]StoredStreamEvent, error)

	// ReadRecords reads events within a single Stream by their names.
	ReadRecords(ctx context.Context, streamName string) ([]StoredStreamEvent, error)
}

type StreamID struct {
	StreamName    string
	StreamVersion uint64
}

type StoredStreamEvent struct {
	ID          StreamID
	EventName   string
	EventData   []byte
	HappenedOn  time.Time
	ContentType string
}
