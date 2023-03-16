package persistence

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type AppendOnlyStore interface {
	// Append appends the marshalled events to the store.
	// It returns an error if the expected version does not match the current version.
	Append(ctx context.Context, name string, data []byte, expectedVersion uint64) error

	// ReadRecords reads events within a single Stream by their names.
	ReadRecords(ctx context.Context, name string, startVersion uint64, maxCount uint64) ([]DataWithVersion, error)

	// ReadAllRecords reads all events from the store. It's commonly used to replicate events.
	ReadAllRecords(ctx context.Context, startVersion uint64, maxCount uint64) ([]DataWithName, error)
}

// DataWithVersion is a data record with a version.
type DataWithVersion struct {
	Version uint64
	Data    []byte
}

type DataWithName struct {
	Name string
	Data []byte
}
