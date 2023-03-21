package account

import (
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func init() {
	serializer.RegisterSerializableType(&AccountOpened{})
	serializer.RegisterSerializableType(&AmountAdded{})
	serializer.RegisterSerializableType(&AmountWithdrawn{})
	serializer.RegisterSerializableType(&TransferenceSent{})
	serializer.RegisterSerializableType(&TransferenceReceived{})
}

type AccountOpened struct {
	domain.BaseEvent

	AccountID ID
}

func NewAccountOpened(id ID) *AccountOpened {
	return &AccountOpened{
		BaseEvent: domain.NewBaseEvent(),
		AccountID: id,
	}
}

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
	From     ID
	To       ID
}

func NewTransferenceSent(quantity int, balance int, from, to ID) *TransferenceSent {
	return &TransferenceSent{
		BaseEvent: domain.NewBaseEvent(),
		Quantity:  quantity,
		Balance:   balance,
		From:      from,
		To:        to,
	}
}

type TransferenceReceived struct {
	domain.BaseEvent

	Quantity int
	Balance  int
	From     ID
	To       ID
}

func NewTransferenceReceived(quantity int, balance int, from, to ID) *TransferenceReceived {
	return &TransferenceReceived{
		BaseEvent: domain.NewBaseEvent(),
		Quantity:  quantity,
		Balance:   balance,
		From:      from,
		To:        to,
	}
}
