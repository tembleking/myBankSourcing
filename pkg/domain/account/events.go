package account

import (
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
	AccountID ID
}

func (a *AccountOpened) EventName() string {
	return "AccountOpened"
}

type AmountAdded struct {
	Quantity int
	Balance  int
}

func (a *AmountAdded) EventName() string {
	return "AmountAdded"
}

type AmountWithdrawn struct {
	Quantity int
	Balance  int
}

func (a *AmountWithdrawn) EventName() string {
	return "AmountWithdrawn"
}

type TransferenceSent struct {
	Quantity int
	Balance  int
	From     ID
	To       ID
}

func (t *TransferenceSent) EventName() string {
	return "TransferenceSent"
}

type TransferenceReceived struct {
	Quantity int
	Balance  int
	From     ID
	To       ID
}

func (t *TransferenceReceived) EventName() string {
	return "TransferenceReceived"
}
