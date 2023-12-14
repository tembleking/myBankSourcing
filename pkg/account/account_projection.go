package account

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Projection struct {
	mutex    sync.RWMutex
	accounts map[string]*Account
}

func (a *Projection) Accounts() []*Account {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	result := make([]*Account, 0, len(a.accounts))

	for _, account := range a.accounts {
		result = append(result, account)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (a *Projection) handleEvent(event domain.Event) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	switch e := event.(type) {
	case *AccountOpened:
		a.accounts[e.AggregateID()] = NewAccount()
		a.accounts[e.AggregateID()].LoadFromHistory(e)
	case *AccountClosed:
		delete(a.accounts, e.AccountID)
	default:
		a.accounts[e.AggregateID()].LoadFromHistory(e)
	}
}

func NewAccountProjection(eventStore *persistence.ReadOnlyEventStore) (*Projection, error) {
	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading all events: %w", err)
	}

	projection := &Projection{accounts: map[string]*Account{}}
	for _, event := range events {
		projection.handleEvent(event)
	}

	return projection, nil
}
