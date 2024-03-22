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

// nolint:revive
type AccountOpened struct {
	Timestamp      time.Time
	ID             domain.EventID
	AccountID      string
	AccountVersion uint64
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
	Timestamp      time.Time
	ID             domain.EventID
	AccountID      string
	Quantity       int
	Balance        int
	AccountVersion uint64
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
	Timestamp      time.Time
	ID             domain.EventID
	AccountID      string
	Quantity       int
	Balance        int
	AccountVersion uint64
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

// nolint:revive
type AccountClosed struct {
	Timestamp      time.Time
	ID             domain.EventID
	AccountID      string
	AccountVersion uint64
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

type TransferSent struct {
	Timestamp          time.Time
	ID                 domain.EventID
	AccountID          string
	TransferID         string
	AccountOrigin      string
	AccountDestination string
	Amount             int
	AccountVersion     uint64
}

func (t *TransferSent) AggregateID() string {
	return t.AccountID
}

func (t *TransferSent) EventID() domain.EventID {
	return t.ID
}

func (t *TransferSent) EventName() string {
	return "TransferSent"
}

func (t *TransferSent) HappenedOn() time.Time {
	return t.Timestamp
}

func (t *TransferSent) Version() uint64 {
	return t.AccountVersion
}

type TransferReceived struct {
	Timestamp          time.Time
	ID                 domain.EventID
	AccountID          string
	TransferID         string
	AccountOrigin      string
	AccountDestination string
	Amount             int
	AccountVersion     uint64
}

func (t *TransferReceived) AggregateID() string {
	return t.AccountID
}

func (t *TransferReceived) EventID() domain.EventID {
	return t.ID
}

func (t *TransferReceived) EventName() string {
	return "TransferReceived"
}

func (t *TransferReceived) HappenedOn() time.Time {
	return t.Timestamp
}

func (t *TransferReceived) Version() uint64 {
	return t.AccountVersion
}

type TransferSentRolledBack struct {
	Timestamp          time.Time
	ID                 domain.EventID
	AccountID          string
	TransferID         string
	AccountOrigin      string
	AccountDestination string
	Amount             int
	AccountVersion     uint64
}

func (t *TransferSentRolledBack) AggregateID() string {
	return t.AccountID
}

func (t *TransferSentRolledBack) EventID() domain.EventID {
	return t.ID
}

func (t *TransferSentRolledBack) EventName() string {
	return "TransferSentRolledBack"
}

func (t *TransferSentRolledBack) HappenedOn() time.Time {
	return t.Timestamp
}

func (t *TransferSentRolledBack) Version() uint64 {
	return t.AccountVersion
}

type TransferCompleted struct {
	Timestamp          time.Time
	ID                 domain.EventID
	AccountID          string
	TransferID         string
	AccountOrigin      string
	AccountDestination string
	Amount             int
	AccountVersion     uint64
}

func (t *TransferCompleted) AggregateID() string {
	return t.AccountID
}

func (t *TransferCompleted) EventID() domain.EventID {
	return t.ID
}

func (t *TransferCompleted) EventName() string {
	return "TransferCompleted"
}

func (t *TransferCompleted) HappenedOn() time.Time {
	return t.Timestamp
}

func (t *TransferCompleted) Version() uint64 {
	return t.AccountVersion
}
