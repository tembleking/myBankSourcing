package transfer

import (
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func init() {
	serializer.RegisterSerializableEvent(&TransferRequested{})
}

// nolint:revive
type TransferRequested struct {
	Timestamp       time.Time
	ID              domain.EventID
	TransferID      string
	FromAccount     string
	ToAccount       string
	Amount          int
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
