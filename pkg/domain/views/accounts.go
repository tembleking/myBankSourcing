package views

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountView struct {
	accountEvents map[account.ID]*account.Account
	rwMutex       sync.RWMutex
}

func (a *AccountView) Dispatch(events ...persistence.StreamEvent) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	for _, event := range events {
		a.handleEvent(event)
	}
}

func (a *AccountView) Accounts() []*account.Account {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	result := make([]*account.Account, 0, len(a.accountEvents))

	for _, account := range a.accountEvents {
		result = append(result, account)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (a *AccountView) handleEvent(event persistence.StreamEvent) {
	switch e := event.Event.(type) {
	case *account.AccountOpened:
		a.accountEvents[e.AccountID] = account.NewAccount(e)
	case *account.AccountClosed:
		delete(a.accountEvents, e.AccountID)
	case *account.AmountAdded:
		a.accountEvents[e.AccountID].Apply(e)
		a.accountEvents[e.AccountID].ClearEvents()
	case *account.AmountWithdrawn:
		a.accountEvents[e.AccountID].Apply(e)
		a.accountEvents[e.AccountID].ClearEvents()
	case *account.TransferSent:
		a.accountEvents[e.From].Apply(e)
		a.accountEvents[e.From].ClearEvents()
	case *account.TransferReceived:
		a.accountEvents[e.To].Apply(e)
		a.accountEvents[e.To].ClearEvents()
	}
}

func NewAccountView(eventStore *persistence.EventStore) (*AccountView, error) {
	a := &AccountView{accountEvents: map[account.ID]*account.Account{}}
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	eventStore.AddDispatchers(a)

	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading all events: %w", err)
	}
	for _, event := range events {
		a.handleEvent(event)
	}

	return a, nil
}
