package transfer

import (
	"context"
	"fmt"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/account"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Transfer struct {
	From     string
	To       string
	Quantity int
}

type TransfersView struct {
	rwMutex   sync.RWMutex
	transfers []Transfer
}

func (v *TransfersView) Transfers() []Transfer {
	v.rwMutex.RLock()
	defer v.rwMutex.RUnlock()
	return v.transfers
}

func (v *TransfersView) handleEvent(event persistence.StreamEvent) {
	if transferSent, ok := event.Event.(*account.TransferSent); ok {
		v.transfers = append(v.transfers, Transfer{
			From:     string(transferSent.From),
			To:       string(transferSent.To),
			Quantity: transferSent.Quantity,
		})
	}
}

func NewTransfersViewFrom(eventStore *persistence.EventStore) (*TransfersView, error) {
	view := &TransfersView{}
	view.rwMutex.Lock()
	defer view.rwMutex.Unlock()

	err := loadViewFromEventStore(eventStore, view)
	if err != nil {
		return nil, fmt.Errorf("error loading view from event store: %w", err)
	}

	return view, nil
}

func loadViewFromEventStore(eventStore *persistence.EventStore, view *TransfersView) error {
	events, err := eventStore.LoadAllEvents(context.Background())
	if err != nil {
		return fmt.Errorf("error loading events: %w", err)
	}
	for _, event := range events {
		view.handleEvent(event)
	}
	return nil
}
