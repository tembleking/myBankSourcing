package account

import (
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Transfer struct {
	TransferID string
	Amount     int
	ToAccount  string
}

func (t *Transfer) ID() string {
	return t.TransferID
}

func (t *Transfer) SameEntityAs(other domain.Entity) bool {
	if otherTransfer, ok := other.(*Transfer); ok {
		return t.ID() == otherTransfer.ID() && t.Amount == otherTransfer.Amount && t.ToAccount == otherTransfer.ToAccount
	}
	return false
}

type Account struct {
	domain.BaseAggregate

	isOpen    bool
	balance   int
	transfers []Transfer
}

func (a *Account) SameEntityAs(other domain.Entity) bool {
	if otherAccount, ok := other.(*Account); ok {
		return a.ID() == otherAccount.ID() &&
			a.Version() == otherAccount.Version() &&
			a.IsOpen() == otherAccount.IsOpen() &&
			a.Balance() == otherAccount.Balance() &&
			slices.EqualFunc(a.transfers, otherAccount.transfers, func(a, b Transfer) bool {
				return a.SameEntityAs(&b)
			})
	}
	return false
}

func NewAccount() *Account {
	a := &Account{}
	a.OnEventFunc = a.onEvent
	return a
}

func OpenAccount(id string) (*Account, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}
	a := NewAccount()
	a.Apply(&AccountOpened{ID: domain.NewEventID(), AccountID: id, AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return a, nil
}

func (a *Account) DepositMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount < 0 {
		return ErrDepositMoneyQuantityCannotBeNegative
	}

	newBalance := a.Balance() + amount
	a.Apply(&AmountDeposited{ID: domain.NewEventID(), AccountID: a.ID(), Quantity: amount, Balance: newBalance, AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return nil
}

func (a *Account) WithdrawMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount > a.Balance() {
		return ErrBalanceIsNotEnoughForWithdrawal
	}

	newBalance := a.Balance() - amount
	a.Apply(&AmountWithdrawn{ID: domain.NewEventID(), AccountID: a.ID(), Quantity: amount, Balance: newBalance, AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return nil
}

func (a *Account) Balance() int {
	return a.balance
}

func (a *Account) onEvent(event domain.Event) {
	switch event := event.(type) {
	case *AccountOpened:
		a.isOpen = true
	case *AmountDeposited:
		a.balance = event.Balance
	case *AmountWithdrawn:
		a.balance = event.Balance
	case *TransferRequested:
		a.balance = event.Balance
		a.transfers = append(a.transfers, Transfer{
			TransferID: event.TransferID,
			Amount:     event.Quantity,
			ToAccount:  event.To,
		})
	case *TransferReceived:
		a.balance = event.Balance
	case *TransferReturned:
		a.balance = event.Balance
	case *AccountClosed:
		a.isOpen = false
	}
}

func (a *Account) IsOpen() bool {
	return a.isOpen
}

func (a *Account) TransferMoney(amount int, destination *Account) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount > a.Balance() {
		return ErrBalanceIsNotEnoughForTransfer
	}

	transferID := uuid.NewString()
	newBalanceOrigin := a.Balance() - amount
	a.Apply(&TransferRequested{ID: domain.NewEventID(), TransferID: transferID, Quantity: amount, Balance: newBalanceOrigin, From: a.ID(), To: destination.ID(), AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return nil
}

func (a *Account) Transfers() []Transfer {
	return a.transfers
}

func (a *Account) CloseAccount() error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if a.Balance() > 0 {
		return ErrAccountCannotBeClosedWithBalance
	}
	a.Apply(&AccountClosed{ID: domain.NewEventID(), AccountID: a.ID(), AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return nil
}

func (a *Account) AcceptTransfer(transferID string, quantity int, originAccount string) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}

	newBalance := a.Balance() + quantity
	a.Apply(&TransferReceived{ID: domain.NewEventID(), TransferID: transferID, Quantity: quantity, Balance: newBalance, From: originAccount, To: a.ID(), AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return nil
}

func (a *Account) ReturnTransfer(transferID string, quantity int, destinationAccount string) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}

	newBalance := a.Balance() + quantity
	a.Apply(&TransferReturned{ID: domain.NewEventID(), TransferID: transferID, Quantity: quantity, Balance: newBalance, From: a.ID(), To: destinationAccount, AccountVersion: a.NextVersion(), Timestamp: a.Now()})
	return nil
}
