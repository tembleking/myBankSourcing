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
		a.OnEventFunc(event)
	}
	return a
}

func (a *Account) ID() ID {
	return a.id
}

func (a *Account) OpenAccount(id ID) error {
	if a.isOpen {
		return ErrAccountIsAlreadyOpen
	}

	a.Apply(NewAccountOpened(id))
	return nil

}

func (a *Account) AddMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount < 0 {
		return ErrAddMoneyQuantityCannotBeNegative
	}

	newBalance := a.Balance() + amount
	a.Apply(NewAmountAdded(amount, newBalance))
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
	a.Apply(NewAmountWithdrawn(amount, newBalance))
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
	case *TransferenceSent:
		a.balance = event.Balance
	case *TransferenceReceived:
		a.balance = event.Balance
	}
}

func (a *Account) IsOpen() bool {
	return a.isOpen
}

func (a *Account) TransferMoney(amount int, destination *Account) error {
	if amount > a.Balance() {
		return ErrBalanceIsNotEnoughForTransfer
	}

	newBalanceOrigin := a.Balance() - amount
	a.Apply(NewTransferenceSent(amount, newBalanceOrigin, a.ID(), destination.ID()))

	newBalanceDestination := destination.Balance() + amount
	destination.Apply(NewTransferenceReceived(amount, newBalanceDestination, a.ID(), destination.ID()))

	return nil
}
