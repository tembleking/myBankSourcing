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

type TransferenceSent struct {
	domain.BaseEvent

	Quantity int
	Balance  int
	To       ID
}

func NewTransferenceSent(quantity int, balance int, to ID) *TransferenceSent {
	return &TransferenceSent{
		BaseEvent: domain.NewBaseEvent(),
		Quantity:  quantity,
		Balance:   balance,
		To:        to,
	}
}

type TransferenceReceived struct {
	domain.BaseEvent

	Quantity int
	Balance  int
	From     ID
}

func NewTransferenceReceived(quantity int, balance int, from ID) *TransferenceReceived {
	return &TransferenceReceived{
		BaseEvent: domain.NewBaseEvent(),
		Quantity:  quantity,
		Balance:   balance,
		From:      from,
	}
}
