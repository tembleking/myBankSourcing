package inmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

type Repository struct {
	accounts map[account.ID][]domain.Event
}

func NewRepository() *Repository {
	return &Repository{
		accounts: map[account.ID][]domain.Event{},
	}
}

func (r *Repository) SaveAccount(ctx context.Context, account *account.Account) error {
	eventsToPersist := account.Events()
	r.accounts[account.ID()] = append(r.accounts[account.ID()], eventsToPersist...)
	return nil
}

func (r *Repository) GetAccount(ctx context.Context, id account.ID) (*account.Account, error) {
	persistedEvents, ok := r.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}

	anAccount := account.NewAccount(id)
	for _, persistedEvent := range persistedEvents {
		// Print is here for demonstration purposes. It can be removed safely.
		fmt.Printf("%#v\n", persistedEvent)
		anAccount.On(persistedEvent)
	}

	return anAccount, nil
}
