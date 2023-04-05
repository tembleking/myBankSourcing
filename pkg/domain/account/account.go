package account

import (
	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type ID string

type Account struct {
	domain.BaseAggregate

	id      ID
	isOpen  bool
	balance int
}

func NewAccount(events ...domain.Event) *Account {
	a := &Account{}
	a.OnEventFunc = a.onEvent
	for _, event := range events {
		a.Apply(event)
	}
	a.ClearEvents()
	return a
}

func (a *Account) ID() ID {
	return a.id
}

func OpenAccount(id ID) *Account {
	a := NewAccount()
	a.Apply(&AccountOpened{AccountID: id, AccountVersion: a.AggregateVersion()})
	return a
}

func (a *Account) AddMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount < 0 {
		return ErrAddMoneyQuantityCannotBeNegative
	}

	newBalance := a.Balance() + amount
	a.Apply(&AmountAdded{AccountID: a.ID(), Quantity: amount, Balance: newBalance, AccountVersion: a.AggregateVersion()})
	return nil
}

func (a *Account) WithdrawalMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount > a.Balance() {
		return ErrBalanceIsNotEnoughForWithdrawal
	}

	newBalance := a.Balance() - amount
	a.Apply(&AmountWithdrawn{AccountID: a.ID(), Quantity: amount, Balance: newBalance, AccountVersion: a.AggregateVersion()})
	return nil
}

func (a *Account) Balance() int {
	return a.balance
}

func (a *Account) onEvent(event domain.Event) {
	switch event := event.(type) {
	case *AccountOpened:
		a.id = event.AccountID
		a.isOpen = true
	case *AmountAdded:
		a.balance = event.Balance
	case *AmountWithdrawn:
		a.balance = event.Balance
	case *TransferSent:
		a.balance = event.Balance
	case *TransferReceived:
		a.balance = event.Balance
	case *AccountClosed:
		a.isOpen = false
	}
}

func (a *Account) IsOpen() bool {
	return a.isOpen
}

func (a *Account) TransferMoney(amount int, destination *Account) error {
	if !a.IsOpen() || !destination.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount > a.Balance() {
		return ErrBalanceIsNotEnoughForTransfer
	}

	newBalanceOrigin := a.Balance() - amount
	a.Apply(&TransferSent{Quantity: amount, Balance: newBalanceOrigin, From: a.ID(), To: destination.ID(), AccountVersion: a.AggregateVersion()})

	newBalanceDestination := destination.Balance() + amount
	destination.Apply(&TransferReceived{Quantity: amount, Balance: newBalanceDestination, From: a.ID(), To: destination.ID(), AccountVersion: destination.AggregateVersion()})

	return nil
}

func (a *Account) CloseAccount() error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if a.Balance() > 0 {
		return ErrAccountCannotBeClosedWithBalance
	}
	a.Apply(&AccountClosed{AccountID: a.ID(), AccountVersion: a.AggregateVersion()})
	return nil
}
