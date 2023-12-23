package account

import (
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/transfer"
)

type Account struct {
	domain.BaseAggregate

	isOpen            bool
	balance           int
	transfersAssigned map[string]struct{}
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
	a := &Account{transfersAssigned: make(map[string]struct{})}
	a.OnEventFunc = a.onEvent
	return a
}

func OpenAccount(id string) (*Account, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
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
	case *TransferAssigned:
		if _, transferAlreadyAssigned := a.transfersAssigned[event.TransferID]; transferAlreadyAssigned {
			return
		}
		if event.AccountOrigin == a.ID() {
			a.balance -= event.Amount
			a.transfersAssigned[event.TransferID] = struct{}{}
		}
		if event.AccountDestination == a.ID() {
			a.balance += event.Amount
			a.transfersAssigned[event.TransferID] = struct{}{}
		}
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

func (a *Account) AssignTransfer(t *transfer.Transfer) error {
	if !a.IsOpen() {
		return ErrAccountIsClosed
	}

	a.Apply(&TransferAssigned{
		ID:                 domain.NewEventID(),
		TransferID:         t.ID(),
		AccountID:          a.ID(),
		AccountOrigin:      t.FromAccount(),
		AccountDestination: t.ToAccount(),
		Amount:             t.Amount(),
		AccountVersion:     a.NextVersion(),
		Timestamp:          a.Now(),
	})
	return nil
}
