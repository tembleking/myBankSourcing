package account

import (
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func init() {
	serializer.RegisterSerializableEvent(&AccountOpened{})
	serializer.RegisterSerializableEvent(&AmountDeposited{})
	serializer.RegisterSerializableEvent(&AmountWithdrawn{})
	serializer.RegisterSerializableEvent(&TransferRequested{})
	serializer.RegisterSerializableEvent(&TransferReceived{})
	serializer.RegisterSerializableEvent(&AccountClosed{})
}

type AccountOpened struct {
	AccountID      string
	AccountVersion uint64
}

func (a *AccountOpened) Version() uint64 {
	return a.AccountVersion
}

func (a *AccountOpened) EventName() string {
	return "AccountOpened"
}

type AmountDeposited struct {
	AccountID      string
	Quantity       int
	Balance        int
	AccountVersion uint64
}

func (a *AmountDeposited) Version() uint64 {
	return a.AccountVersion
}

func (a *AmountDeposited) EventName() string {
	return "AmountAdded"
}

type AmountWithdrawn struct {
	AccountID      string
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

type TransferRequested struct {
	TransferID     string
	Quantity       int
	Balance        int
	From           string
	To             string
	AccountVersion uint64
}

func (t *TransferRequested) EventName() string {
	return "TransferRequested"
}

func (t *TransferRequested) Version() uint64 {
	return t.AccountVersion
}

type TransferReceived struct {
	Quantity       int
	Balance        int
	From           string
	To             string
	AccountVersion uint64
}

func (t *TransferReceived) EventName() string {
	return "TransferReceived"
}

func (t *TransferReceived) Version() uint64 {
	return t.AccountVersion
}

type AccountClosed struct {
	AccountID      string
	AccountVersion uint64
}

func (a *AccountClosed) EventName() string {
	return "AccountClosed"
}

func (a *AccountClosed) Version() uint64 {
	return a.AccountVersion
}
