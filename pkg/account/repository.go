package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

var ErrAccountNotFound = errors.New("account not found")

type Repository struct {
	eventStore *persistence.EventStore
}

func (r *Repository) NextID() string {
	return uuid.NewString()
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Account, error) {
	events, err := r.eventStore.LoadEventStream(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve events from event store: %w", err)
	}

	account := NewAccount()
	account.LoadFromHistory(events...)
	return account, err
}

func (r *Repository) Save(ctx context.Context, aggregate *Account) error {
	err := r.eventStore.AppendToStream(ctx, aggregate)
	if err != nil {
		return fmt.Errorf("unable to append to event store: %w", err)
	}
	return nil
}

func NewRepository(eventStore *persistence.EventStore) *Repository {
	return &Repository{
		eventStore: eventStore,
	}
}
