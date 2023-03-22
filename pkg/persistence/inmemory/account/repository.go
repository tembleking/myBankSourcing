package account

import (
	"context"
	"errors"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

type Repository struct {
	accounts map[account.ID][]domain.Event

	outbox        chan domain.Event
	rwMutex       sync.RWMutex
	subscriptions map[chan domain.Event]struct{}
}

func (r *Repository) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	accounts := make([]*account.Account, 0)
	for _, events := range r.accounts {
		accounts = append(accounts, account.NewAccount(events...))
	}

	return accounts, nil
}

func NewRepository() *Repository {
	r := &Repository{
		accounts:      map[account.ID][]domain.Event{},
		outbox:        make(chan domain.Event),
		subscriptions: map[chan domain.Event]struct{}{},
	}
	return r
}

func (r *Repository) SaveAccount(ctx context.Context, account *account.Account) error {
	eventsToPersist := account.Events()
	r.accounts[account.ID()] = append(r.accounts[account.ID()], eventsToPersist...)
	go r.sendEventsToSubscriptions(eventsToPersist)
	return nil
}

func (r *Repository) GetAccount(ctx context.Context, id account.ID) (*account.Account, error) {
	persistedEvents, ok := r.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}

	anAccount := account.NewAccount(persistedEvents...)
	return anAccount, nil
}

func (r *Repository) Subscribe(ctx context.Context) (<-chan domain.Event, error) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	subscription := make(chan domain.Event)
	r.subscriptions[subscription] = struct{}{}

	go func() {
		<-ctx.Done()
		r.rwMutex.Lock()
		defer r.rwMutex.Unlock()
		delete(r.subscriptions, subscription)
		close(subscription)
	}()

	return subscription, nil
}

func (r *Repository) sendEventsToSubscriptions(events []domain.Event) {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	for subscription := range r.subscriptions {
		for _, event := range events {
			subscription <- event
		}
	}
}
