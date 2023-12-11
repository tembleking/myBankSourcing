package account

import (
	"time"

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
	Timestamp      time.Time
}

func (a *AccountOpened) AggregateID() string {
	return a.AccountID
}

func (a *AccountOpened) Version() uint64 {
	return a.AccountVersion
}

func (a *AccountOpened) EventName() string {
	return "AccountOpened"
}

func (a *AccountOpened) HappenedOn() time.Time {
	return a.Timestamp
}

type AmountDeposited struct {
	AccountID      string
	Quantity       int
	Balance        int
	AccountVersion uint64
	Timestamp      time.Time
}

func (a *AmountDeposited) AggregateID() string {
	return a.AccountID
}

func (a *AmountDeposited) Version() uint64 {
	return a.AccountVersion
}

func (a *AmountDeposited) EventName() string {
	return "AmountDeposited"
}

func (a *AmountDeposited) HappenedOn() time.Time {
	return a.Timestamp
}

type AmountWithdrawn struct {
	AccountID      string
	Quantity       int
	Balance        int
	AccountVersion uint64
	Timestamp      time.Time
}

func (a *AmountWithdrawn) AggregateID() string {
	return a.AccountID
}

func (a *AmountWithdrawn) Version() uint64 {
	return a.AccountVersion
}

func (a *AmountWithdrawn) EventName() string {
	return "AmountWithdrawn"
}

func (a *AmountWithdrawn) HappenedOn() time.Time {
	return a.Timestamp
}

type TransferRequested struct {
	TransferID     string
	Quantity       int
	Balance        int
	From           string
	To             string
	AccountVersion uint64
	Timestamp      time.Time
}

func (t *TransferRequested) AggregateID() string {
	return t.From
}

func (t *TransferRequested) EventName() string {
	return "TransferRequested"
}

func (t *TransferRequested) Version() uint64 {
	return t.AccountVersion
}

func (a *TransferRequested) HappenedOn() time.Time {
	return a.Timestamp
}

type TransferReceived struct {
	TransferID     string
	Quantity       int
	Balance        int
	From           string
	To             string
	AccountVersion uint64
	Timestamp      time.Time
}

func (t *TransferReceived) AggregateID() string {
	return t.To
}

func (t *TransferReceived) EventName() string {
	return "TransferReceived"
}

func (t *TransferReceived) Version() uint64 {
	return t.AccountVersion
}

func (a *TransferReceived) HappenedOn() time.Time {
	return a.Timestamp
}

type AccountClosed struct {
	AccountID      string
	AccountVersion uint64
	Timestamp      time.Time
}

func (a *AccountClosed) AggregateID() string {
	return a.AccountID
}

func (a *AccountClosed) EventName() string {
	return "AccountClosed"
}

func (a *AccountClosed) Version() uint64 {
	return a.AccountVersion
}

func (a *AccountClosed) HappenedOn() time.Time {
	return a.Timestamp
}
