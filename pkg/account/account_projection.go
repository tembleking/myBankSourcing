package account

import (
	"context"
	"fmt"
	"sort"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountProjection struct {
	accountEvents map[string]*Account
}

func (a *AccountProjection) Accounts() []*Account {
	result := make([]*Account, 0, len(a.accountEvents))

	for _, account := range a.accountEvents {
		result = append(result, account)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (a *AccountProjection) handleEvent(event domain.Event) {
	switch e := event.(type) {
	case *AccountOpened:
		a.accountEvents[e.AccountID] = NewAccount(e)
	case *AccountClosed:
		delete(a.accountEvents, e.AccountID)
	case *AmountDeposited:
		a.accountEvents[e.AccountID].LoadFromHistory(e)
	case *AmountWithdrawn:
		a.accountEvents[e.AccountID].LoadFromHistory(e)
	case *TransferRequested:
		a.accountEvents[e.From].LoadFromHistory(e)
	case *TransferReceived:
		a.accountEvents[e.To].LoadFromHistory(e)
	}
}

func NewAccountProjection(eventStore *persistence.EventStore) (*AccountProjection, error) {
	a := &AccountProjection{accountEvents: map[string]*Account{}}

	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading all events: %w", err)
	}
	for _, event := range events {
		a.handleEvent(event.Event)
	}

	return a, nil
}
