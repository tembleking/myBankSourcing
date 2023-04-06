package views

import (
	"context"
	"sort"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountView struct {
	accountEvents map[account.ID][]domain.Event
	rwMutex       sync.RWMutex
}

func (a *AccountView) Dispatch(ctx context.Context, events ...persistence.StreamEvent) {
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

	for _, events := range a.accountEvents {
		result = append(result, account.NewAccount(events...))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (a *AccountView) handleEvent(event persistence.StreamEvent) {
	switch e := event.Event.(type) {
	case *account.AccountOpened:
		a.accountEvents[e.AccountID] = append(a.accountEvents[e.AccountID], e)
	case *account.AccountClosed:
		delete(a.accountEvents, e.AccountID)
	case *account.AmountAdded:
		a.accountEvents[e.AccountID] = append(a.accountEvents[e.AccountID], e)
	case *account.AmountWithdrawn:
		a.accountEvents[e.AccountID] = append(a.accountEvents[e.AccountID], e)
	case *account.TransferSent:
		a.accountEvents[e.From] = append(a.accountEvents[e.From], e)
	case *account.TransferReceived:
		a.accountEvents[e.To] = append(a.accountEvents[e.To], e)
	}
}

func NewAccountView(eventStore *persistence.EventStore) (*AccountView, error) {
	a := &AccountView{accountEvents: map[account.ID][]domain.Event{}}
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	eventStore.AddDispatchers(a)

	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		a.handleEvent(event)
	}

	return a, nil
}
