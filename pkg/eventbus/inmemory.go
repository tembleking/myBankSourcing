package eventbus

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type InMemory struct {
	listeners []domain.EventListener
	mutex     sync.RWMutex
}

func (i *InMemory) Publish(ctx context.Context, events ...domain.Event) error {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, event := range events {
		if err := i.sendEventToListeners(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (i *InMemory) sendEventToListeners(ctx context.Context, event domain.Event) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, listener := range i.listeners {
		listener := listener
		group.Go(func() error {
			return listener.OnEvent(ctx, event)
		})
	}
	return group.Wait()
}

func (i *InMemory) Subscribe(_ context.Context, listener domain.EventListener) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.listeners = append(i.listeners, listener)

	return nil
}

func NewInMemory() domain.EventBus {
	return &InMemory{}
}
