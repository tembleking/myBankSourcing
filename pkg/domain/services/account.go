package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/domain/views"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountService struct {
	eventStore  *persistence.EventStore
	accountView *views.AccountView
}

func (a *AccountService) OpenAccount(ctx context.Context) (*account.Account, error) {
	accountCreated := account.OpenAccount(account.ID(uuid.NewString()))

	err := a.eventStore.AppendToStream(ctx, string(accountCreated.ID()), accountCreated.AggregateVersion(), accountCreated.Events())
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return accountCreated, nil
}

func (a *AccountService) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	return a.accountView.Accounts(), nil
}

func (a *AccountService) AddMoneyToAccount(ctx context.Context, accountID account.ID, amount int) (*account.Account, error) {
	account, err := a.getAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.AddMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error adding money to account: %w", err)
	}

	err = a.eventStore.AppendToStream(ctx, string(account.ID()), account.AggregateVersion(), account.Events())
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return account, nil
}

func (a *AccountService) WithdrawMoneyFromAccount(ctx context.Context, accountID account.ID, amount int) (*account.Account, error) {
	account, err := a.getAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.WithdrawMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error withdrawing money from account: %w", err)
	}

	err = a.eventStore.AppendToStream(ctx, string(account.ID()), account.AggregateVersion(), account.Events())
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return account, nil
}

func (a *AccountService) CloseAccount(ctx context.Context, accountID account.ID) (*account.Account, error) {
	account, err := a.getAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.CloseAccount()
	if err != nil {
		return nil, fmt.Errorf("error closing account: %w", err)
	}

	err = a.eventStore.AppendToStream(ctx, string(account.ID()), account.AggregateVersion(), account.Events())
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return account, nil
}

func (a *AccountService) TransferMoney(ctx context.Context, origin account.ID, destination account.ID, amountToTransfer int) (*account.Account, error) {
	originAccount, err := a.getAccount(ctx, origin)
	if err != nil {
		return nil, fmt.Errorf("error getting origin account: %w", err)
	}

	destinationAccount, err := a.getAccount(ctx, destination)
	if err != nil {
		return nil, fmt.Errorf("error getting destination account: %w", err)
	}

	err = originAccount.TransferMoney(amountToTransfer, destinationAccount)
	if err != nil {
		return nil, fmt.Errorf("error transferring money: %w", err)
	}

	err = a.eventStore.AppendToStream(ctx, string(originAccount.ID()), originAccount.AggregateVersion(), originAccount.Events())
	if err != nil {
		return nil, fmt.Errorf("error saving from account: %w", err)
	}

	err = a.eventStore.AppendToStream(ctx, string(destinationAccount.ID()), destinationAccount.AggregateVersion(), destinationAccount.Events())
	if err != nil {
		return nil, fmt.Errorf("error saving to account: %w", err)
	}

	return originAccount, nil
}

func (a *AccountService) getAccount(ctx context.Context, id account.ID) (*account.Account, error) {
	stream, err := a.eventStore.LoadEventStream(ctx, string(id))
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

func NewAccountService(eventStore *persistence.EventStore) *AccountService {
	accountView, err := views.NewAccountView(eventStore)
	if err != nil {
		panic(err)
	}
	return &AccountService{
		eventStore:  eventStore,
		accountView: accountView,
	}
}
