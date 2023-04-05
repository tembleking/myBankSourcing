package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

type AccountService struct {
	repository account.Repository
}

func (s *AccountService) OpenAccount(ctx context.Context) (*account.Account, error) {
	accountCreated := account.OpenAccount(account.ID(uuid.NewString()))
	err := s.repository.SaveAccount(ctx, accountCreated)
	if err != nil {
		return nil, fmt.Errorf("error saving account: %w", err)
	}

	return accountCreated, nil
}

func NewAccountService(repository account.Repository) *AccountService {
	return &AccountService{
		repository: repository,
	}
}
