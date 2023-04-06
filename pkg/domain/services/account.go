package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

type AccountRepository interface {
	SaveAccount(ctx context.Context, account *account.Account) error
	GetAccount(ctx context.Context, accountID account.ID) (*account.Account, error)
	ListAccounts(ctx context.Context) ([]*account.Account, error)
}

type AccountService struct {
	repository AccountRepository
}

func (s *AccountService) OpenAccount(ctx context.Context) (*account.Account, error) {
	accountCreated := account.OpenAccount(account.ID(uuid.NewString()))
	err := s.repository.SaveAccount(ctx, accountCreated)
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return accountCreated, nil
}

func (s *AccountService) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	return s.repository.ListAccounts(ctx)
}

func (s *AccountService) AddMoneyToAccount(ctx context.Context, accountID account.ID, amount int) (*account.Account, error) {
	account, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.AddMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error adding money to account: %w", err)
	}

	err = s.repository.SaveAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return account, nil
}

func NewAccountService(accountRepository AccountRepository) *AccountService {
	return &AccountService{
		repository: accountRepository,
	}
}
