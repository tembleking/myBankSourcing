package account

import (
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func init() {
	serializer.RegisterSerializableType(&AccountOpened{})
	serializer.RegisterSerializableType(&AmountAdded{})
	serializer.RegisterSerializableType(&AmountWithdrawn{})
	serializer.RegisterSerializableType(&TransferSent{})
	serializer.RegisterSerializableType(&TransferReceived{})
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
	AccountID      ID
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
	AccountID      ID
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

type TransferSent struct {
	Quantity       int
	Balance        int
	From           ID
	To             ID
	AccountVersion uint64
}

func (t *TransferSent) EventName() string {
	return "TransferSent"
}

func (t *TransferSent) Version() uint64 {
	return t.AccountVersion
}

type TransferReceived struct {
	Quantity       int
	Balance        int
	From           ID
	To             ID
	AccountVersion uint64
}

func (t *TransferReceived) EventName() string {
	return "TransferReceived"
}

func (t *TransferReceived) Version() uint64 {
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
