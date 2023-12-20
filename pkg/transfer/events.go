package transfer

import (
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type TransferRequested struct {
	ID              domain.EventID
	TransferID      string
	FromAccount     string
	ToAccount       string
	Amount          int
	Timestamp       time.Time
	TransferVersion uint64
}

func (t *TransferRequested) AggregateID() string {
	return t.TransferID
}

func (t *TransferRequested) EventID() domain.EventID {
	return t.ID
}

func (t *TransferRequested) EventName() string {
	return "TransferRequested"
}

func (t *TransferRequested) HappenedOn() time.Time {
	return t.Timestamp
}

func (t *TransferRequested) Version() uint64 {
	return t.TransferVersion
}
