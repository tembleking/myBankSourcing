package transferences

import (
	"context"
	"sync"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

type View struct {
	rwMutex       sync.RWMutex
	transferences []Transference
}

func (v *View) listenAndUpdateFromEvents(ctx context.Context, subscription <-chan domain.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-subscription:
			if !ok {
				return
			}
			v.rwMutex.Lock()
			v.onEvent(event)
			v.rwMutex.Unlock()
		}
	}
}

func (v *View) onEvent(event domain.Event) {
	switch event := event.(type) {
	case *account.TransferenceSent:
		v.transferences = append(v.transferences, Transference{
			Origin:      event.From,
			Destination: event.To,
			Quantity:    event.Quantity,
		})
	}
}

type Transference struct {
	Origin      account.ID
	Destination account.ID
	Quantity    int
}

func (v *View) Transferences() ([]Transference, error) {
	v.rwMutex.RLock()
	defer v.rwMutex.RUnlock()

	return v.transferences, nil
}

func NewViewSubscribedTo(ctx context.Context, subscribable domain.Subscribable) (*View, error) {
	subscription, err := subscribable.Subscribe(ctx)
	if err != nil {
		return nil, err
	}

	v := &View{}
	go v.listenAndUpdateFromEvents(ctx, subscription)

	return v, nil
}
