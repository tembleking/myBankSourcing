package services

import (
	"context"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

type OpenAccountService struct {
	repository account.Repository
}

func (s *OpenAccountService) OpenAccount(background context.Context) (*account.Account, error) {
	accountCreated := account.NewAccount()

	err := s.repository.SaveAccount(background, accountCreated)
	if err != nil {
		return nil, err
	}

	return accountCreated, nil
}

func NewCreateAccountService(repository account.Repository) *OpenAccountService {
	return &OpenAccountService{
		repository: repository,
	}
}
