package account

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Transfer struct {
	TransferID string
	Amount     int
	ToAccount  string
}

type Account struct {
	domain.BaseAggregate

	id        string
	isOpen    bool
	balance   int
	transfers []Transfer
}

func NewAccount(events ...domain.Event) *Account {
	a := &Account{}
	a.OnEventFunc = a.onEvent
	for _, event := range events {
		err := a.Apply(event)
		if err != nil {
			panic(err)
		}
	}
	a.ClearEvents()
	return a
}

func (a *Account) ID() string {
	return a.id
}

func OpenAccount(id string) (*Account, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}
	a := NewAccount()
	err := a.Apply(&AccountOpened{AccountID: id, AccountVersion: a.Version()})
	return a, err
}

func (a *Account) DepositMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount < 0 {
		return ErrDepositMoneyQuantityCannotBeNegative
	}

	newBalance := a.Balance() + amount
	return a.Apply(&AmountDeposited{AccountID: a.ID(), Quantity: amount, Balance: newBalance, AccountVersion: a.Version()})
}

func (a *Account) WithdrawMoney(amount int) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}
	if amount > a.Balance() {
		return ErrBalanceIsNotEnoughForWithdrawal
	}

	newBalance := a.Balance() - amount
	return a.Apply(&AmountWithdrawn{AccountID: a.ID(), Quantity: amount, Balance: newBalance, AccountVersion: a.Version()})
}

func (a *Account) Balance() int {
	return a.balance
}

func (a *Account) onEvent(event domain.Event) {
	switch event := event.(type) {
	case *AccountOpened:
		a.id = event.AccountID
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
	if err := a.Apply(&TransferRequested{TransferID: transferID, Quantity: amount, Balance: newBalanceOrigin, From: a.ID(), To: destination.ID(), AccountVersion: a.Version()}); err != nil {
		return err
	}

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
	return a.Apply(&AccountClosed{AccountID: a.ID(), AccountVersion: a.Version()})
}
