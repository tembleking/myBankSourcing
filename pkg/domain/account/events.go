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
	AccountID      ID
	AccountVersion uint64
}

func (a *AccountOpened) Version() uint64 {
	return a.AccountVersion
}

func (a *AccountOpened) EventName() string {
	return "AccountOpened"
}

type AmountAdded struct {
	Quantity       int
	Balance        int
	AccountVersion uint64
}

func (a *AmountAdded) Version() uint64 {
	return a.AccountVersion
}

func (a *AmountAdded) EventName() string {
	return "AmountAdded"
}

type AmountWithdrawn struct {
	Quantity       int
	Balance        int
	AccountVersion uint64
}

func (a *AmountWithdrawn) Version() uint64 {
	return a.AccountVersion
}

func (a *AmountWithdrawn) EventName() string {
	return "AmountWithdrawn"
}

type TransferenceSent struct {
	Quantity       int
	Balance        int
	From           ID
	To             ID
	AccountVersion uint64
}

func (t *TransferenceSent) EventName() string {
	return "TransferenceSent"
}

func (t *TransferenceSent) Version() uint64 {
	return t.AccountVersion
}

type TransferenceReceived struct {
	Quantity       int
	Balance        int
	From           ID
	To             ID
	AccountVersion uint64
}

func (t *TransferenceReceived) EventName() string {
	return "TransferenceReceived"
}

func (t *TransferenceReceived) Version() uint64 {
	return t.AccountVersion
}

type AccountClosed struct {
	AccountID      ID
	AccountVersion uint64
}

func (a *AccountClosed) EventName() string {
	return "AccountClosed"
}

func (a *AccountClosed) Version() uint64 {
	return a.AccountVersion
}
