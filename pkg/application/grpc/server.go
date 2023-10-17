package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/application/proto"
)

type AccountGRPCServer struct {
	accountService *account.AccountService
	accountView    *account.AccountProjection
}

func NewAccountGRPCServer(accountService *account.AccountService, accountView *account.AccountProjection) *AccountGRPCServer {
	return &AccountGRPCServer{
		accountService: accountService,
		accountView:    accountView,
	}
}

func (s *AccountGRPCServer) OpenAccount(ctx context.Context, empty *emptypb.Empty) (*proto.OpenAccountResponse, error) {
	account, err := s.accountService.OpenAccount(ctx)
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 500, Err: err}
	}

	return &proto.OpenAccountResponse{
		Account: &proto.Account{
			Id:      string(account.ID()),
			Balance: int64(account.Balance()),
		},
	}, nil
}

func (s *AccountGRPCServer) ListAccounts(ctx context.Context, empty *emptypb.Empty) (*proto.ListAccountsResponse, error) {
	accounts := s.accountView.Accounts()
	protoAccounts := make([]*proto.Account, len(accounts))
	for i, account := range accounts {
		protoAccounts[i] = &proto.Account{
			Id:      string(account.ID()),
			Balance: int64(account.Balance()),
		}
	}
	return &proto.ListAccountsResponse{
		Accounts: protoAccounts,
	}, nil
}

func (s *AccountGRPCServer) AddMoney(ctx context.Context, request *proto.AddMoneyRequest) (*proto.AddMoneyResponse, error) {
	accountID := request.GetAccountId()
	if accountID == "" {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 400, Err: fmt.Errorf("account id must be provided")}
	}
	amount := int(request.GetAmount())
	if amount <= 0 {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 400, Err: fmt.Errorf("amount must be greater than 0")}
	}

	account, err := s.accountService.AddMoneyToAccount(ctx, accountID, amount)
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 500, Err: err}
	}

	return &proto.AddMoneyResponse{
		Account: &proto.Account{
			Id:      string(account.ID()),
			Balance: int64(account.Balance()),
		},
	}, nil

}

func (s *AccountGRPCServer) WithdrawMoney(ctx context.Context, request *proto.WithdrawMoneyRequest) (*proto.WithdrawMoneyResponse, error) {
	accountID := request.GetAccountId()
	amount := int(request.GetAmount())
	account, err := s.accountService.WithdrawMoneyFromAccount(ctx, accountID, amount)
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 500, Err: err}
	}

	return &proto.WithdrawMoneyResponse{
		Account: &proto.Account{
			Id:      string(account.ID()),
			Balance: int64(account.Balance()),
		},
	}, nil
}

func (s *AccountGRPCServer) TransferMoney(ctx context.Context, request *proto.TransferMoneyRequest) (*proto.TransferMoneyResponse, error) {
	fromAccountID := request.GetFromAccountId()
	toAccountID := request.GetToAccountId()
	amount := int(request.GetAmount())
	account, err := s.accountService.TransferMoney(ctx, fromAccountID, toAccountID, amount)
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 500, Err: err}
	}

	return &proto.TransferMoneyResponse{
		Account: &proto.Account{
			Id:      string(account.ID()),
			Balance: int64(account.Balance()),
		},
	}, nil
}

func (s *AccountGRPCServer) CloseAccount(ctx context.Context, request *proto.CloseAccountRequest) (*emptypb.Empty, error) {
	accountID := request.GetAccountId()
	_, err := s.accountService.CloseAccount(ctx, accountID)
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: 500, Err: err}
	}

	return &emptypb.Empty{}, nil
}
