package account

import (
	"context"
	"fmt"
	"sort"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/domain/views"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Repository struct {
	eventStore  *persistence.EventStore
	accountView *views.AccountView
}

func NewRepository(store *persistence.EventStore) *Repository {
	accountView, err := views.NewAccountView(store)
	if err != nil {
		panic(err)
	}

	return &Repository{
		eventStore:  store,
		accountView: accountView,
	}
}

func (r *Repository) SaveAccount(ctx context.Context, account *account.Account) error {
	err := r.eventStore.AppendToStream(ctx, string(account.ID()), account.AggregateVersion(), account.Events())
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
	if len(stream) == 0 {
		return nil, fmt.Errorf("account not found: %s", id)
	}

	events := make([]domain.Event, 0, len(stream))
	for _, event := range stream {
		events = append(events, event.Event)
	}

	return account.NewAccount(events...), nil
}

func (r *Repository) ListAccounts(ctx context.Context) ([]account.ID, error) {
	allAccounts := r.accountView.Accounts()

	accountIDs := make([]account.ID, 0, len(allAccounts))
	for _, account := range allAccounts {
		accountIDs = append(accountIDs, account.ID())
	}

	sort.Slice(accountIDs, func(i, j int) bool {
		return accountIDs[i] < accountIDs[j]
	})

	return accountIDs, nil
}
