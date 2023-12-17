package account

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type Projection struct {
	mutex                 sync.RWMutex
	accounts              map[string]*Account
	precalculatedAccounts []Account
	lastProcessedEventID  domain.EventID
	eventStore            *persistence.ReadOnlyEventStore
}

func (a *Projection) Accounts() []Account {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.precalculatedAccounts
}

func (a *Projection) handleEvent(event domain.Event) {
	switch e := event.(type) {
	case *AccountOpened:
		a.accounts[e.AggregateID()] = NewAccount()
		a.accounts[e.AggregateID()].LoadFromHistory(e)
	case *AccountClosed:
		delete(a.accounts, e.AccountID)
	default:
		a.accounts[e.AggregateID()].LoadFromHistory(e)
	}

	a.lastProcessedEventID = event.EventID()
}

func (a *Projection) precalculateAccounts() {
	accounts := make([]Account, 0, len(a.accounts))

	for _, account := range a.accounts {
		accounts = append(accounts, *account)
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].ID() < accounts[j].ID()
	})

	a.precalculatedAccounts = accounts
}

func (a *Projection) refreshProjection(ctx context.Context) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	store := a.eventStore
	if a.lastProcessedEventID != "" {
		store = store.AfterEventID(a.lastProcessedEventID)
	}

	events, err := store.LoadAllEvents(ctx)
	if err != nil {
		slog.Default().ErrorContext(ctx, "error loading events from store", "error", err.Error())
		return
	}

	for _, event := range events {
		a.handleEvent(event)
	}

	if len(events) > 0 {
		a.precalculateAccounts()
	}
}

func (a *Projection) startPeriodicRefresh(ctx context.Context, refreshInterval time.Duration) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.refreshProjection(ctx)
		}
	}
}

func NewAccountProjection(ctx context.Context, eventStore *persistence.ReadOnlyEventStore, refreshInterval time.Duration) (*Projection, error) {
	p := &Projection{accounts: map[string]*Account{}, eventStore: eventStore}
	p.refreshProjection(ctx)
	go p.startPeriodicRefresh(ctx, refreshInterval)
	return p, nil
}
