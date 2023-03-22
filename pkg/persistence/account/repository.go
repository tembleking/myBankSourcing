package account

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Repository struct {
	eventStore *persistence.EventStore
}

func NewRepository(store *persistence.EventStore) *Repository {
	return &Repository{
		eventStore: store,
	}
}

func (r *Repository) SaveAccount(ctx context.Context, account *account.Account) error {
	eventsToPersist := account.Events()
	expectedVersion := uint64(0)
	stream, err := r.eventStore.LoadEventStream(ctx, string(account.ID()))
	if err == nil {
		expectedVersion = stream.Version
	}

	err = r.eventStore.AppendToStream(ctx, string(account.ID()), expectedVersion, eventsToPersist)
	if err != nil {
		return fmt.Errorf("error appending to stream: %w", err)
	}

	return nil
}

func (r *Repository) GetAccount(ctx context.Context, id account.ID) (*account.Account, error) {
	stream, err := r.eventStore.LoadEventStream(ctx, string(id))
	if err != nil {
		return nil, fmt.Errorf("error loading event stream: %w", err)
	}

	return account.NewAccount(stream.Events...), nil
}
