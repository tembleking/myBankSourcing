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

type ProjectedAccount struct {
	AccountID string
	Movements []ProjectedMovement
	Balance   int
}

type ProjectedMovement struct {
	Timestamp        time.Time
	Type             string
	Amount           int
	ResultingBalance int
}

type Projection struct {
	accounts              map[string]*ProjectedAccount
	eventStore            *persistence.ReadOnlyEventStore
	lastProcessedEventID  domain.EventID
	precalculatedAccounts []ProjectedAccount
	mutex                 sync.RWMutex
}

func (a *Projection) Accounts() []ProjectedAccount {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.precalculatedAccounts
}

func (a *Projection) handleEvent(event domain.Event) {
	switch e := event.(type) {
	case *AccountOpened:
		a.accounts[e.AggregateID()] = &ProjectedAccount{AccountID: e.AggregateID()}
	case *AccountClosed:
		delete(a.accounts, e.AccountID)
	case *AmountDeposited:
		a.accounts[e.AggregateID()].Balance += e.Quantity
		a.accounts[e.AggregateID()].Movements = append(a.accounts[e.AggregateID()].Movements, ProjectedMovement{
			Type:             "Deposit",
			Amount:           e.Quantity,
			ResultingBalance: e.Balance,
			Timestamp:        e.HappenedOn(),
		})
	case *AmountWithdrawn:
		a.accounts[e.AggregateID()].Balance -= e.Quantity
		a.accounts[e.AggregateID()].Movements = append(a.accounts[e.AggregateID()].Movements, ProjectedMovement{
			Type:             "Withdrawal",
			Amount:           e.Quantity,
			ResultingBalance: e.Balance,
			Timestamp:        e.HappenedOn(),
		})
	}

	a.lastProcessedEventID = event.EventID()
}

func (a *Projection) precalculateAccounts() {
	accounts := make([]ProjectedAccount, 0, len(a.accounts))

	for _, account := range a.accounts {
		accounts = append(accounts, *account)
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].AccountID < accounts[j].AccountID
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
	p := &Projection{accounts: map[string]*ProjectedAccount{}, eventStore: eventStore}
	p.refreshProjection(ctx)
	go p.startPeriodicRefresh(ctx, refreshInterval)
	return p, nil
}
