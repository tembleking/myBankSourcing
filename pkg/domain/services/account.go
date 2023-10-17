package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountService struct {
	eventStore        *persistence.EventStore
	accountRepository domain.Repository[*account.Account]
}

func (a *AccountService) OpenAccount(ctx context.Context) (*account.Account, error) {
	accountCreated := account.OpenAccount(uuid.NewString())

	err := a.accountRepository.Save(ctx, accountCreated)
	if err != nil {
		return nil, fmt.Errorf("error saving created account: %w", err)
	}

	return accountCreated, err
}

func (a *AccountService) AddMoneyToAccount(ctx context.Context, accountID string, amount int) (*account.Account, error) {
	account, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.AddMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error adding money to account: %w", err)
	}

	err = a.accountRepository.Save(ctx, account)

	return account, err
}

func (a *AccountService) WithdrawMoneyFromAccount(ctx context.Context, accountID string, amount int) (*account.Account, error) {
	account, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.WithdrawMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error withdrawing money from account: %w", err)
	}

	err = a.accountRepository.Save(ctx, account)

	return account, err
}

func (a *AccountService) CloseAccount(ctx context.Context, accountID string) (*account.Account, error) {
	account, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.CloseAccount()
	if err != nil {
		return nil, fmt.Errorf("error closing account: %w", err)
	}

	err = a.accountRepository.Save(ctx, account)

	return account, err
}

func (a *AccountService) TransferMoney(ctx context.Context, origin string, destination string, amountToTransfer int) (*account.Account, error) {
	originAccount, err := a.accountRepository.GetByID(ctx, origin)
	if err != nil {
		return nil, fmt.Errorf("error getting origin account: %w", err)
	}

	destinationAccount, err := a.accountRepository.GetByID(ctx, destination)
	if err != nil {
		return nil, fmt.Errorf("error getting destination account: %w", err)
	}

	err = originAccount.TransferMoney(amountToTransfer, destinationAccount)
	if err != nil {
		return nil, fmt.Errorf("error transferring money: %w", err)
	}

	err = a.eventStore.AppendToStream(ctx, originAccount, destinationAccount)
	if err != nil {
		return nil, fmt.Errorf("error saving from accounts: %w", err)
	}

	return originAccount, err
}

func NewAccountService(eventStore *persistence.EventStore, accountRepository domain.Repository[*account.Account]) *AccountService {
	return &AccountService{
		eventStore:        eventStore,
		accountRepository: accountRepository,
	}
}
