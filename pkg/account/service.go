package account

import (
	"context"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/transfer"
)

type Service struct {
	accountRepository  domain.Repository[*Account]
	transferRepository domain.Repository[*transfer.Transfer]
}

func (a *Service) OnCommand(ctx context.Context, command domain.Command) error {
	if c, ok := command.(*OpenNewAccount); ok {
		_, err := a.openAccount(ctx, c.ID)
		return err
	}

	return nil
}

func (a *Service) OpenAccount(ctx context.Context) (*Account, error) {
	return a.openAccount(ctx, a.accountRepository.NextID())
}

func (a *Service) openAccount(ctx context.Context, accountID string) (*Account, error) {
	accountCreated, err := OpenAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("error opening account: %w", err)
	}

	if err := a.accountRepository.Save(ctx, accountCreated); err != nil {
		return nil, fmt.Errorf("error saving created account: %w", err)
	}

	return accountCreated, err
}

func (a *Service) DepositMoneyIntoAccount(ctx context.Context, accountID string, amount int) (*Account, error) {
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

func (a *Service) WithdrawMoneyFromAccount(ctx context.Context, accountID string, amount int) (*Account, error) {
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

func (a *Service) CloseAccount(ctx context.Context, accountID string) (*Account, error) {
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

func (a *Service) TransferMoney(ctx context.Context, originAccountID string, destinationAccountID string, amount int) (*transfer.Transfer, error) {
	origin, err := a.accountRepository.GetByID(ctx, originAccountID)
	if err != nil {
		return nil, fmt.Errorf("error getting origin account: %w", err)
	}

	destination, err := a.accountRepository.GetByID(ctx, destinationAccountID)
	if err != nil {
		return nil, fmt.Errorf("error getting destination account: %w", err)
	}

	transfer, err := origin.TransferMoney(amount, destination)
	if err != nil {
		return nil, fmt.Errorf("error creating transfer: %w", err)
	}

	err = a.transferRepository.Save(ctx, transfer)
	if err != nil {
		return nil, fmt.Errorf("error saving transfer: %w", err)
	}

	return transfer, nil
}

func (a *Service) SendTransfer(ctx context.Context, transferID string) error {
	transfer, err := a.transferRepository.GetByID(ctx, transferID)
	if err != nil {
		return fmt.Errorf("error getting the transfer: %w", err)
	}

	originAccount, err := a.accountRepository.GetByID(ctx, transfer.FromAccount())
	if err != nil {
		return fmt.Errorf("error getting the origin account: %w", err)
	}

	err = originAccount.SendTransfer(transfer)
	if err != nil {
		return fmt.Errorf("error sending the transfer: %w", err)
	}

	err = a.accountRepository.Save(ctx, originAccount)
	if err != nil {
		return fmt.Errorf("error saving the account: %w", err)
	}

	return nil
}

func (a *Service) ReceiveTransfer(ctx context.Context, transferID string) error {
	transfer, err := a.transferRepository.GetByID(ctx, transferID)
	if err != nil {
		return fmt.Errorf("error getting the transfer: %w", err)
	}

	destinationAccount, err := a.accountRepository.GetByID(ctx, transfer.ToAccount())
	if err != nil {
		return fmt.Errorf("error getting the destination account: %w", err)
	}

	err = destinationAccount.ReceiveTransfer(transfer)
	if err != nil {
		return fmt.Errorf("error sending the transfer: %w", err)
	}

	err = a.accountRepository.Save(ctx, destinationAccount)
	if err != nil {
		return fmt.Errorf("error saving the account: %w", err)
	}

	return nil
}

func (a *Service) RollbackTransfer(ctx context.Context, transferID string) error {
	transfer, err := a.transferRepository.GetByID(ctx, transferID)
	if err != nil {
		return fmt.Errorf("error getting the transfer: %w", err)
	}

	originAccount, err := a.accountRepository.GetByID(ctx, transfer.FromAccount())
	if err != nil {
		return fmt.Errorf("error getting the origin account: %w", err)
	}

	err = originAccount.RollbackSentTransfer(transfer)
	if err != nil {
		return fmt.Errorf("error sending the transfer: %w", err)
	}

	err = a.accountRepository.Save(ctx, originAccount)
	if err != nil {
		return fmt.Errorf("error saving the account: %w", err)
	}

	return nil
}

func (a *Service) CompleteTransfer(ctx context.Context, transferID string) error {
	transfer, err := a.transferRepository.GetByID(ctx, transferID)
	if err != nil {
		return fmt.Errorf("error getting the transfer: %w", err)
	}

	originAccount, err := a.accountRepository.GetByID(ctx, transfer.FromAccount())
	if err != nil {
		return fmt.Errorf("error getting the origin account: %w", err)
	}

	err = originAccount.MarkTransferAsCompleted(transfer)
	if err != nil {
		return fmt.Errorf("error sending the transfer: %w", err)
	}

	err = a.accountRepository.Save(ctx, originAccount)
	if err != nil {
		return fmt.Errorf("error saving the account: %w", err)
	}

	return nil
}

func NewAccountService(accountRepository domain.Repository[*Account], transferRepository domain.Repository[*transfer.Transfer]) *Service {
	return &Service{
		accountRepository:  accountRepository,
		transferRepository: transferRepository,
	}
}
