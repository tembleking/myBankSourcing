package account

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountService struct {
	eventStore        *persistence.EventStore
	accountRepository domain.Repository[*Account]
}

func (a *AccountService) OpenAccount(ctx context.Context) (*Account, error) {
	accountCreated, err := OpenAccount(a.accountRepository.NextID())
	if err != nil {
		return nil, fmt.Errorf("error opening account: %w", err)
	}

	if err := a.accountRepository.Save(ctx, accountCreated); err != nil {
		return nil, fmt.Errorf("error saving created account: %w", err)
	}

	return accountCreated, err
}

func (a *AccountService) DepositMoneyIntoAccount(ctx context.Context, accountID string, amount int) (*Account, error) {
	account, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.DepositMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error depositing money to account: %w", err)
	}

	err = a.accountRepository.Save(ctx, account)

	return account, err
}

func (a *AccountService) WithdrawMoneyFromAccount(ctx context.Context, accountID string, amount int) (*Account, error) {
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

func (a *AccountService) CloseAccount(ctx context.Context, accountID string) (*Account, error) {
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

func NewAccountService(eventStore *persistence.EventStore, accountRepository domain.Repository[*Account]) *AccountService {
	return &AccountService{
		eventStore:        eventStore,
		accountRepository: accountRepository,
	}
}