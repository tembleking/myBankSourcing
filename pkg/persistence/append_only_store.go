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

	ReadAllRecords(ctx context.Context) ([]StoredStreamEvent, error)

	// ReadRecords reads events within a single Stream by their names.
	ReadRecords(ctx context.Context, streamID string) ([]StoredStreamEvent, error)

	ReadEventsByName(ctx context.Context, eventName string) ([]StoredStreamEvent, error)
}

type StoredStreamEvent struct {
	StreamID      string
	StreamVersion uint64
	EventName     string
	EventData     []byte
	HappenedOn    time.Time
}
