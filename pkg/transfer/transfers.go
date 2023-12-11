package transfer

import (
	"context"
	"fmt"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Transfer struct {
	ID       string
	From     string
	To       string
	Quantity int
}

type TransfersProjection struct {
	rwMutex   sync.RWMutex
	transfers []Transfer
}

func (v *TransfersProjection) Transfers() []Transfer {
	v.rwMutex.RLock()
	defer v.rwMutex.RUnlock()
	return v.transfers
}

func (v *TransfersProjection) handleEvent(event domain.Event) {
	if transferSent, ok := event.(*account.TransferRequested); ok {
		v.transfers = append(v.transfers, Transfer{
			ID:       transferSent.TransferID,
			From:     string(transferSent.From),
			To:       string(transferSent.To),
			Quantity: transferSent.Quantity,
		})
	}
}

func NewTransfersProjectionFrom(eventStore *persistence.EventStore) (*TransfersProjection, error) {
	projection := &TransfersProjection{}
	projection.rwMutex.Lock()
	defer projection.rwMutex.Unlock()

	err := loadProjectionFromEventStore(eventStore, projection)
	if err != nil {
		return nil, fmt.Errorf("error loading projection from event store: %w", err)
	}

	return projection, nil
}

func loadProjectionFromEventStore(eventStore *persistence.EventStore, projection *TransfersProjection) error {
	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return fmt.Errorf("error loading events: %w", err)
	}
	for _, event := range events {
		projection.handleEvent(event)
	}
	return nil
}
