package transfer_test

import (
	"context"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/transfer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

var _ = Describe("Transfers", func() {
	var (
		eventStore *persistence.EventStore
	)

	BeforeEach(func() {
		eventStore = persistence.NewEventStoreBuilder(sqlite.InMemory()).Build()
		anAggregate := fakeAggregate{}.withID("some-account").withVersion(5).withEvents([]domain.Event{
			&account.AccountOpened{AccountID: "some-account", AccountVersion: 0},
			&account.AmountDeposited{Quantity: 50, Balance: 50, AccountVersion: 1},
			&account.AmountDeposited{Quantity: 100, Balance: 150, AccountVersion: 2},
			&account.TransferRequested{Quantity: 50, Balance: 100, From: "some-account", To: "another-account", AccountVersion: 3},
			&account.AmountWithdrawn{Quantity: 25, Balance: 75, AccountVersion: 4},
		}...)

		err := eventStore.AppendToStream(context.Background(), &anAggregate)
		Expect(err).ToNot(HaveOccurred())
	})

	When("there are multiple TransferSent events saved", func() {
		It("returns the transfers found", func() {
			transfersProjection, err := transfer.NewTransfersProjectionFrom(eventStore)
			Expect(err).ToNot(HaveOccurred())

			transfers := transfersProjection.Transfers()

			Expect(transfers).To(ConsistOf(
				transfer.Transfer{From: "some-account", To: "another-account", Quantity: 50},
			))
		})
	})
})

type fakeAggregate struct {
	domain.BaseAggregate
	id      string
	version uint64
	events  []domain.Event
}

// ID implements domain.Aggregate.
func (f *fakeAggregate) ID() string {
	return f.id
}

// Version implements domain.Aggregate.
func (f *fakeAggregate) Version() uint64 {
	return f.version
}

func (f *fakeAggregate) SameEntityAs(other domain.Entity) bool {
	//TODO implement me
	panic("implement me")
}

func (f fakeAggregate) withID(id string) fakeAggregate {
	return fakeAggregate{
		id:      id,
		version: f.version,
		events:  f.events,
	}
}

func (f fakeAggregate) withVersion(version uint64) fakeAggregate {
	return fakeAggregate{
		id:      f.id,
		version: version,
		events:  f.events,
	}
}

func (f fakeAggregate) withEvents(events ...domain.Event) fakeAggregate {
	aggregate := fakeAggregate{
		id:      f.id,
		version: f.version,
		events:  events,
	}
	for _, event := range events {
		aggregate.Apply(event)
	}
	return aggregate
}
