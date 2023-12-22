package account

import (
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func init() {
	serializer.RegisterSerializableEvent(&AccountOpened{})
	serializer.RegisterSerializableEvent(&AmountDeposited{})
	serializer.RegisterSerializableEvent(&AmountWithdrawn{})
	serializer.RegisterSerializableEvent(&AccountClosed{})
}

type AccountOpened struct {
	ID             domain.EventID
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

func (a *AccountOpened) EventID() domain.EventID {
	return a.ID
}

func (a *AccountOpened) EventName() string {
	return "AccountOpened"
}

func (a *AccountOpened) HappenedOn() time.Time {
	return a.Timestamp
}

type AmountDeposited struct {
	ID             domain.EventID
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

func (a *AmountDeposited) EventID() domain.EventID {
	return a.ID
}

func (a *AmountDeposited) EventName() string {
	return "AmountDeposited"
}

func (a *AmountDeposited) HappenedOn() time.Time {
	return a.Timestamp
}

type AmountWithdrawn struct {
	ID             domain.EventID
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

func (a *AmountWithdrawn) EventID() domain.EventID {
	return a.ID
}

func (a *AmountWithdrawn) EventName() string {
	return "AmountWithdrawn"
}

func (a *AmountWithdrawn) HappenedOn() time.Time {
	return a.Timestamp
}

type AccountClosed struct {
	ID             domain.EventID
	AccountID      string
	AccountVersion uint64
	Timestamp      time.Time
}

func (a *AccountClosed) AggregateID() string {
	return a.AccountID
}

func (a *AccountClosed) EventID() domain.EventID {
	return a.ID
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

type TransferAssigned struct {
	ID                 domain.EventID
	AccountID          string
	TransferID         string
	AccountOrigin      string
	AccountDestination string
	Amount             int
	AccountVersion     uint64
	Timestamp          time.Time
}

func (t *TransferAssigned) AggregateID() string {
	return t.AccountID
}

func (t *TransferAssigned) EventID() domain.EventID {
	return t.ID
}

func (t *TransferAssigned) EventName() string {
	return "TransferAssigned"
}

func (t *TransferAssigned) HappenedOn() time.Time {
	return t.Timestamp
}

func (t *TransferAssigned) Version() uint64 {
	return t.AccountVersion
}
