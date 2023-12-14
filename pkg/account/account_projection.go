package account

import (
	"context"
	"fmt"
	"sort"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Projection struct {
	accountEvents map[string]*Account
}

func (a *Projection) Accounts() []*Account {
	result := make([]*Account, 0, len(a.accountEvents))

	for _, account := range a.accountEvents {
		result = append(result, account)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (a *Projection) handleEvent(event domain.Event) {
	switch e := event.(type) {
	case *AccountOpened:
		a.accountEvents[e.AggregateID()] = NewAccount()
		a.accountEvents[e.AggregateID()].LoadFromHistory(e)
	case *AccountClosed:
		delete(a.accountEvents, e.AccountID)
	case *AmountDeposited, *AmountWithdrawn, *TransferRequested, *TransferReceived:
		a.accountEvents[e.AggregateID()].LoadFromHistory(e)
	}
}

func NewAccountProjection(eventStore *persistence.ReadOnlyEventStore) (*Projection, error) {
	a := &Projection{accountEvents: map[string]*Account{}}

	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading all events: %w", err)
	}
	for _, event := range events {
		a.handleEvent(event)
	}

	return a, nil
}
