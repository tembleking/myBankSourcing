package account

import "github.com/tembleking/myBankSourcing/pkg/domain"

type ID string

type Account struct {
	domain.BaseAggregate

	id      ID
	balance int
}

func NewAccount(id ID) *Account {
	a := &Account{id: id}
	a.OnEventFunc = a.onEvent
	return a
}

func (a *Account) ID() ID {
	return a.id
}

func (a *Account) AddMoney(amount int) error {
	if amount < 0 {
		return ErrAddMoneyQuantityCannotBeNegative
	}

	newBalance := a.Balance() + amount
	a.Apply(NewAmountAdded(amount, newBalance))
	return nil
}

func (a *Account) WithdrawalMoney(amount int) error {
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
	case *AmountAdded:
		a.balance = event.Balance
	case *AmountWithdrawn:
		a.balance = event.Balance
	}
}
