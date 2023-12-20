package account

import (
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Account struct {
	domain.BaseAggregate

	isOpen  bool
	balance int
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
	case *AccountClosed:
		a.isOpen = false
	}
}

func (a *Account) IsOpen() bool {
	return a.isOpen
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
