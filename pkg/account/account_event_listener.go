package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type EventListener struct {
	accountRepository domain.Repository[*Account]
}

func (e *EventListener) OnEvent(ctx context.Context, event domain.Event) error {
	switch event := event.(type) {
	case *TransferRequested:
		return e.handleEventTransferRequested(ctx, event)
	}
	return nil
}

func (e *EventListener) handleEventTransferRequested(ctx context.Context, event *TransferRequested) error {
	destinationAccount, err := e.accountRepository.GetByID(ctx, event.To)
	if errors.Is(err, ErrAccountNotFound) {
		return e.returnTransferMoneyToOriginAccount(ctx, event)
	}
	if err != nil {
		return fmt.Errorf("error getting destination account: %w", err)
	}

	if err := destinationAccount.AcceptTransfer(event.TransferID, event.Quantity, event.From); err != nil {
		return e.returnTransferMoneyToOriginAccount(ctx, event)
	}

	if err = e.accountRepository.Save(ctx, destinationAccount); err != nil {
		return fmt.Errorf("error saving transfer destination account: %w", err)
	}

	return nil
}

func (e *EventListener) returnTransferMoneyToOriginAccount(ctx context.Context, event *TransferRequested) error {
	originAccount, err := e.accountRepository.GetByID(ctx, event.From)
	if err != nil {
		return fmt.Errorf("error getting origin account: %w", err)
	}

	if err := originAccount.ReturnTransfer(event.TransferID, event.Quantity, event.To); err != nil {
		return fmt.Errorf("error returning transfer to origin account: %w", err)
	}

	if err = e.accountRepository.Save(ctx, originAccount); err != nil {
		return fmt.Errorf("error saving transfer origin account: %w", err)
	}

	return nil
}

func NewEventListener(accountRepository domain.Repository[*Account]) *EventListener {
	return &EventListener{
		accountRepository: accountRepository,
	}
}
