package account

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Repository struct {
	eventStore *persistence.EventStore
}

// GetByID implements domain.Repository.
func (r *Repository) GetByID(ctx context.Context, id string) (*Account, error) {
	events, err := r.eventStore.LoadEventStream(ctx, persistence.StreamName(id))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve events from event store: %w", err)
	}

	domainEvents := make([]domain.Event, 0, len(events))
	for _, event := range events {
		domainEvents = append(domainEvents, event.Event)
	}

	return NewAccount(domainEvents...), err
}

// Save implements domain.Repository.
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
