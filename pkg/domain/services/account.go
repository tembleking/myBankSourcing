package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AccountRepository interface {
	SaveAccount(ctx context.Context, account *account.Account) error
	GetAccount(ctx context.Context, accountID account.ID) (*account.Account, error)
	ListAccounts(ctx context.Context) ([]*account.Account, error)
}

type AccountService struct {
	repository AccountRepository
	eventStore *persistence.EventStore
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

func (s *AccountService) WithdrawMoneyFromAccount(ctx context.Context, accountID account.ID, amount int) (*account.Account, error) {
	account, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.WithdrawMoney(amount)
	if err != nil {
		return nil, fmt.Errorf("error withdrawing money from account: %w", err)
	}

	err = s.repository.SaveAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return account, nil
}

func (s *AccountService) CloseAccount(ctx context.Context, accountID account.ID) (*account.Account, error) {
	account, err := s.repository.GetAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	err = account.CloseAccount()
	if err != nil {
		return nil, fmt.Errorf("error closing account: %w", err)
	}

	err = s.repository.SaveAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return account, nil
}

func (s *AccountService) TransferMoney(ctx context.Context, origin account.ID, destination account.ID, amountToTransfer int) (*account.Account, error) {
	originAccount, err := s.repository.GetAccount(ctx, origin)
	if err != nil {
		return nil, fmt.Errorf("error getting origin account: %w", err)
	}

	destinationAccount, err := s.repository.GetAccount(ctx, destination)
	if err != nil {
		return nil, fmt.Errorf("error getting destination account: %w", err)
	}

	err = originAccount.TransferMoney(amountToTransfer, destinationAccount)
	if err != nil {
		return nil, fmt.Errorf("error transferring money: %w", err)
	}

	err = s.repository.SaveAccount(ctx, originAccount)
	if err != nil {
		return nil, fmt.Errorf("error saving from account: %w", err)
	}

	err = s.repository.SaveAccount(ctx, destinationAccount)
	if err != nil {
		return nil, fmt.Errorf("error saving to account: %w", err)
	}

	return originAccount, nil
}

func NewAccountService(accountRepository AccountRepository, eventStore *persistence.EventStore) *AccountService {
	return &AccountService{
		repository: accountRepository,
		eventStore: eventStore,
	}
}
