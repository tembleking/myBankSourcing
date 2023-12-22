package transfer

import (
	"context"
	"errors"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

var ErrTransferNotFound = errors.New("transfer not found")

type Repository struct {
	eventStore *persistence.EventStore
}

func (r *Repository) NextID() string {
	return domain.NewUUID()
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Transfer, error) {
	events, err := r.eventStore.LoadEventStream(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve events from event store: %w", err)
	}

	account := NewTransfer()
	account.LoadFromHistory(events...)
	return account, err
}

func (r *Repository) Save(ctx context.Context, aggregate *Transfer) error {
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
