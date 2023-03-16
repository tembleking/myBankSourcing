package inmemory

import (
	"context"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type dataWithVersionAndName struct {
	data    []byte
	version uint64
	name    string
}

type AppendOnlyStore struct {
	fields  map[string][]dataWithVersionAndName
	rwMutex sync.RWMutex
}

func (a *AppendOnlyStore) Append(ctx context.Context, name string, data []byte, expectedVersion uint64) error {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	existingVersion := 0
	fields, ok := a.fields[name]
	if !ok {
		fields = make([]dataWithVersionAndName, 0)
	}

	if len(fields) > 0 {
		existingVersion = len(fields)
	}

	if uint64(existingVersion) != expectedVersion {
		return &persistence.ErrUnexpectedVersion{Found: uint64(existingVersion), Expected: expectedVersion}
	}

	a.fields[name] = append(fields, dataWithVersionAndName{
		version: uint64(existingVersion + 1),
		data:    data,
		name:    name,
	})

	return nil
}

func (a *AppendOnlyStore) ReadRecords(ctx context.Context, name string, startVersion uint64, maxCount uint64) ([]persistence.DataWithVersion, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	fields, ok := a.fields[name]
	if !ok {
		return nil, &persistence.ErrAggregateNotFound{Name: name}
	}

	if startVersion > uint64(len(fields)) {
		return nil, nil
	}

	if maxCount == 0 || maxCount > uint64(len(fields)) {
		maxCount = uint64(len(fields))
	}

	result := make([]persistence.DataWithVersion, 0)
	for i := startVersion; i < uint64(len(fields)) && maxCount > 0; i++ {
		maxCount--
		result = append(result, persistence.DataWithVersion{
			Data:    fields[i].data,
			Version: fields[i].version,
		})
	}
	return result, nil
}

func (a *AppendOnlyStore) ReadAllRecords(ctx context.Context, startVersion uint64, maxCount uint64) ([]persistence.DataWithName, error) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	result := make([]persistence.DataWithName, 0)
	for _, fields := range a.fields {
		records, err := a.ReadRecords(ctx, fields[0].name, startVersion, maxCount)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			result = append(result, persistence.DataWithName{
				Data: record.Data,
				Name: fields[0].name,
			})
		}
	}
	return result, nil
}

func NewAppendOnlyStore() *AppendOnlyStore {
	return &AppendOnlyStore{
		fields: make(map[string][]dataWithVersionAndName),
	}
}
