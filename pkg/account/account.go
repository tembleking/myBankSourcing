package account

import (
	"errors"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/transfer"
)

type Account struct {
	transfersSent                map[string]struct{}
	transfersReceived            map[string]struct{}
	transfersRolledBack          map[string]struct{}
	pendingTransfersToBeResolved map[string]struct{}
	domain.BaseAggregate

	balance int
	isOpen  bool
}

func (a *Account) SameEntityAs(other domain.Entity) bool {
	if otherAccount, ok := other.(*Account); ok {
		return a.ID() == otherAccount.ID() &&
			a.Version() == otherAccount.Version() &&
			a.IsOpen() == otherAccount.IsOpen() &&
			a.Balance() == otherAccount.Balance()
	}
	return false
}

func NewAccount() *Account {
	a := &Account{
		transfersSent:                make(map[string]struct{}),
		transfersReceived:            make(map[string]struct{}),
		transfersRolledBack:          make(map[string]struct{}),
		pendingTransfersToBeResolved: make(map[string]struct{}),
	}
	a.OnEventFunc = a.onEvent
	return a
}

func OpenAccount(id string) (*Account, error) {
	if id == "" {
		return nil, errors.New("id must not be empty")
	}
	a := NewAccount()
	a.Apply(&AccountOpened{
		ID:             domain.NewEventID(),
		AccountID:      id,
		AccountVersion: a.NextVersion(),
		Timestamp:      a.Now(),
	})
	return a, nil
}

func (a *Account) DepositMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount < 0 {
		return ErrQuantityCannotBeNegative
	}

	newBalance := a.Balance() + amount
	a.Apply(&AmountDeposited{
		ID:             domain.NewEventID(),
		AccountID:      a.ID(),
		Quantity:       amount,
		Balance:        newBalance,
		AccountVersion: a.NextVersion(),
		Timestamp:      a.Now(),
	})
	return nil
}

func (a *Account) WithdrawMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount > a.Balance() {
		return ErrBalanceIsNotEnough
	}
	if amount < 0 {
		return ErrQuantityCannotBeNegative
	}

	newBalance := a.Balance() - amount
	a.Apply(&AmountWithdrawn{
		ID:             domain.NewEventID(),
		AccountID:      a.ID(),
		Quantity:       amount,
		Balance:        newBalance,
		AccountVersion: a.NextVersion(),
		Timestamp:      a.Now(),
	})
	return nil
}

func (a *Account) CloseAccount() error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if a.Balance() > 0 {
		return ErrAccountCannotBeClosedWithBalance
	}
	if len(a.pendingTransfersToBeResolved) > 0 {
		return ErrAccountCannotBeClosedUntilTransfersAreResolved
	}

	a.Apply(&AccountClosed{
		ID:             domain.NewEventID(),
		AccountID:      a.ID(),
		AccountVersion: a.NextVersion(),
		Timestamp:      a.Now(),
	})
	return nil
}

func (a *Account) TransferMoney(amount int, destination *Account) (*transfer.Transfer, error) {
	if a.ID() == destination.ID() {
		return nil, ErrCannotTransferToSameAccount
	}
	if !a.IsOpen() || !destination.IsOpen() {
		return nil, ErrAccountIsClosed
	}
	if a.Balance() < amount {
		return nil, ErrBalanceIsNotEnough
	}
	if amount < 0 {
		return nil, ErrQuantityCannotBeNegative
	}

	return transfer.RequestTransfer(a.ID(), destination.ID(), amount), nil
}

func (a *Account) SendTransfer(transfer *transfer.Transfer) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if a.isTransferAlreadySent(transfer) {
		return nil // idempotent
	}

	a.Apply(&TransferSent{
		ID:                 domain.NewEventID(),
		TransferID:         transfer.ID(),
		AccountID:          a.ID(),
		AccountOrigin:      transfer.FromAccount(),
		AccountDestination: transfer.ToAccount(),
		Amount:             transfer.Amount(),
		AccountVersion:     a.NextVersion(),
		Timestamp:          a.Now(),
	})
	return nil
}

func (a *Account) ReceiveTransfer(transfer *transfer.Transfer) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if a.isTransferAlreadyReceived(transfer) {
		return nil // idempotent
	}

	a.Apply(&TransferReceived{
		ID:                 domain.NewEventID(),
		TransferID:         transfer.ID(),
		AccountID:          a.ID(),
		AccountOrigin:      transfer.FromAccount(),
		AccountDestination: transfer.ToAccount(),
		Amount:             transfer.Amount(),
		AccountVersion:     a.NextVersion(),
		Timestamp:          a.Now(),
	})
	return nil
}

func (a *Account) RollbackSentTransfer(transfer *transfer.Transfer) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}

	if !a.isTransferAlreadySent(transfer) {
		return ErrCannotRollbackTransferNotPreviouslySent
	}

	if a.isTransferAlreadyRolledBack(transfer) {
		return nil // idempotent
	}

	a.Apply(&TransferSentRolledBack{
		ID:                 domain.NewEventID(),
		TransferID:         transfer.ID(),
		AccountID:          a.ID(),
		AccountOrigin:      transfer.FromAccount(),
		AccountDestination: transfer.ToAccount(),
		Amount:             transfer.Amount(),
		AccountVersion:     a.NextVersion(),
		Timestamp:          a.Now(),
	})

	return nil
}

func (a *Account) MarkTransferAsCompleted(transfer *transfer.Transfer) error {
	if !a.isTransferAlreadySent(transfer) {
		return ErrCannotCompleteTransferNotPreviouslySent
	}
	if !a.isTransferPendingToBeResolved(transfer) {
		return nil // idempotent
	}

	a.Apply(&TransferCompleted{
		ID:                 domain.NewEventID(),
		TransferID:         transfer.ID(),
		AccountID:          a.ID(),
		AccountOrigin:      transfer.FromAccount(),
		AccountDestination: transfer.ToAccount(),
		Amount:             transfer.Amount(),
		AccountVersion:     a.NextVersion(),
		Timestamp:          a.Now(),
	})
	return nil
}

func (a *Account) Balance() int {
	return a.balance
}

func (a *Account) IsOpen() bool {
	return a.isOpen
}

func (a *Account) isTransferAlreadySent(transfer *transfer.Transfer) bool {
	_, transferAlreadySent := a.transfersSent[transfer.ID()]
	return transferAlreadySent
}

func (a *Account) isTransferAlreadyReceived(transfer *transfer.Transfer) bool {
	_, transferAlreadyReceived := a.transfersReceived[transfer.ID()]
	return transferAlreadyReceived
}

func (a *Account) isTransferAlreadyRolledBack(transfer *transfer.Transfer) bool {
	_, transferAlreadyRolledBack := a.transfersRolledBack[transfer.ID()]
	return transferAlreadyRolledBack
}

func (a *Account) isTransferPendingToBeResolved(transfer *transfer.Transfer) bool {
	_, transferPendingToBeResolved := a.pendingTransfersToBeResolved[transfer.ID()]
	return transferPendingToBeResolved
}

func (a *Account) onEvent(event domain.Event) {
	switch event := event.(type) {
	case *AccountOpened:
		a.isOpen = true
	case *AmountDeposited:
		a.balance = event.Balance
	case *AmountWithdrawn:
		a.balance = event.Balance
	case *AccountClosed:
		a.isOpen = false
	case *TransferSent:
		a.balance -= event.Amount
		a.transfersSent[event.TransferID] = struct{}{}
		a.pendingTransfersToBeResolved[event.TransferID] = struct{}{}
	case *TransferReceived:
		a.balance += event.Amount
		a.transfersReceived[event.TransferID] = struct{}{}
	case *TransferSentRolledBack:
		a.balance += event.Amount
		a.transfersRolledBack[event.TransferID] = struct{}{}
		delete(a.pendingTransfersToBeResolved, event.TransferID)
	case *TransferCompleted:
		delete(a.pendingTransfersToBeResolved, event.TransferID)

	}
}
