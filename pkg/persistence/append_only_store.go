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
	ReadRecords(ctx context.Context, streamID string) ([]StoredStreamEvent, error)

	// ReadEventsByName reads events by their names.
	ReadEventsByName(ctx context.Context, eventName string) ([]StoredStreamEvent, error)

	// ReadUndispatchedRecords reads all events that have not been dispatched.
	// It should not return the same events twice, unless they have not been marked as dispatched after some time.
	ReadUndispatchedRecords(ctx context.Context) ([]StoredStreamEvent, error)

	// MarkRecordsAsDispatched marks the events as dispatched.
	MarkRecordsAsDispatched(ctx context.Context, events ...StoredStreamEvent) error
}

type StoredStreamEvent struct {
	StreamID      string
	StreamVersion uint64
	EventName     string
	EventData     []byte
	HappenedOn    time.Time
}

func (s StoredStreamEvent) Equal(other StoredStreamEvent) bool {
	return s.StreamID == other.StreamID &&
		s.StreamVersion == other.StreamVersion &&
		s.EventName == other.EventName &&
		byteSliceEqual(s.EventData, other.EventData) &&
		s.HappenedOn.Equal(other.HappenedOn)
}

func byteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
