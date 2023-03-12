package account

import (
	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type AmountAdded struct {
	domain.BaseEvent

	Quantity int
	Balance  int
}

func NewAmountAdded(quantity int, balance int) *AmountAdded {
	return &AmountAdded{
		BaseEvent: domain.NewBaseEvent(),
		Quantity:  quantity,
		Balance:   balance,
	}
}

type AmountWithdrawn struct {
	domain.BaseEvent

	Quantity int
	Balance  int
}

func NewAmountWithdrawn(quantity int, balance int) *AmountWithdrawn {
	return &AmountWithdrawn{
		BaseEvent: domain.NewBaseEvent(),
		Quantity:  quantity,
		Balance:   balance,
	}
}
