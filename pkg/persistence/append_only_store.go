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
	ReadRecords(ctx context.Context, streamName StreamName) ([]StoredStreamEvent, error)

	// ReadEventsByName reads events by their names.
	ReadEventsByName(ctx context.Context, eventName string) ([]StoredStreamEvent, error)

	// ReadUndispatchedRecords reads all events that have not been dispatched.
	// It should not return the same events twice, unless they have not been marked as dispatched after some time.
	ReadUndispatchedRecords(ctx context.Context) ([]StoredStreamEvent, error)

	// MarkRecordsAsDispatched marks the event as dispatched.
	MarkRecordsAsDispatched(ctx context.Context, events ...StreamID) error
}

type StoredStreamEvent struct {
	ID         StreamID
	EventName  string
	EventData  []byte
	HappenedOn time.Time
}
