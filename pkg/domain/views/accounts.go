package views

import (
	"context"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountView struct {
	accounts []string
}

func (v *AccountView) Accounts() []string {
	return v.accounts
}

func (v *AccountView) handleEvent(event persistence.StreamEvent) {
	switch e := event.Event.(type) {
	case *account.AccountOpened:
		v.accounts = append(v.accounts, string(e.AccountID))
	}
}

func NewAccountView(eventStore *persistence.EventStore) (*AccountView, error) {
	a := &AccountView{}

	events, err := eventStore.LoadEventsByName(context.Background(), "AccountOpened")
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		a.handleEvent(event)
	}

	return a, nil
}
