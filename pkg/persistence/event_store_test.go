package persistence_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/mocks"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
	. "github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("EventStore", func() {
	var (
		ctx             context.Context
		ctrl            *gomock.Controller
		appendOnlyStore *mocks.MockAppendOnlyStore
		eventStore      *persistence.EventStore
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		appendOnlyStore = mocks.NewMockAppendOnlyStore(ctrl)
		eventStore = persistence.NewEventStoreBuilder(sqlite.InMemory()).WithAppendOnlyStore(appendOnlyStore).Build()
	})

	It("should be able to load an event stream", func() {
		appendOnlyStore.EXPECT().ReadRecords(ctx, "aggregate-0").Return(
			[]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventID: "event0", EventData: dataRecordInStore(), EventName: "AmountDeposited"}},
			nil,
		)

		stream, err := eventStore.LoadEventStream(ctx, "aggregate-0")

		Expect(err).To(BeNil())
		Expect(stream).To(Equal([]domain.Event{
			&account.AmountDeposited{ID: "event0", AccountID: "some-account", Quantity: 10, Balance: 10},
		}))
	})

	It("should be able to append to an event stream", func() {
		appendOnlyStore.EXPECT().Append(
			ctx,
			persistence.StoredStreamEvent{
				ID:          persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0},
				EventID:     "event0",
				EventName:   "AmountDeposited",
				EventData:   dataRecordInStore(),
				ContentType: "application/json",
			},
			persistence.StoredStreamEvent{
				ID:          persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0},
				EventID:     "event0",
				EventName:   "AmountDeposited",
				EventData:   dataRecordInStore(),
				ContentType: "application/json",
			},
		).Return(nil)

		anAggregate := fakeAggregate{}.withID("aggregate-0").withVersion(1).withEvents(
			&account.AmountDeposited{ID: "event0", AccountID: "some-account", Quantity: 10, Balance: 10},
		)
		anotherAggregate := fakeAggregate{}.withID("aggregate-1").withVersion(1).withEvents(
			&account.AmountDeposited{ID: "event0", AccountID: "some-account", Quantity: 10, Balance: 10},
		)
		err := eventStore.AppendToStream(ctx, &anAggregate, &anotherAggregate)
		Expect(err).To(BeNil())
	})

	When("asking for all events", func() {
		It("returns all events", func() {
			appendOnlyStore.EXPECT().
				ReadAllRecords(ctx).
				Return([]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventID: "event0", EventData: dataRecordInStore(), EventName: "AmountDeposited"}}, nil)

			stream, err := eventStore.LoadAllEvents(ctx)
			Expect(err).To(BeNil())
			Expect(stream).To(Equal([]domain.Event{
				&account.AmountDeposited{ID: "event0", AccountID: "some-account", Quantity: 10, Balance: 10},
			}))
		})
	})

	When("asking for the events after some event ID", func() {
		It("returns a new store with these modifications", func() {
			appendOnlyStore.EXPECT().AfterEventID("event0").Return(appendOnlyStore)

			newEventStore := eventStore.AfterEventID("event0")
			Expect(newEventStore).NotTo(BeNil())
			Expect(newEventStore).NotTo(PointToTheSameLocationAs(eventStore))
		})
	})

	When("asking for a limit number of events", func() {
		It("returns a new store with these modifications", func() {
			appendOnlyStore.EXPECT().Limit(42).Return(appendOnlyStore)

			newEventStore := eventStore.Limit(42)
			Expect(newEventStore).NotTo(BeNil())
			Expect(newEventStore).NotTo(PointToTheSameLocationAs(eventStore))
		})
	})
})

func dataRecordInStore() []byte {
	serializer := &serializer.JSON{}
	data, err := serializer.SerializeDomainEvent(
		&account.AmountDeposited{ID: "event0", AccountID: "some-account", Quantity: 10, Balance: 10},
	)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return data
}

type fakeAggregate struct {
	domain.BaseAggregate
	id      string
	version uint64
}

func (f *fakeAggregate) SameEntityAs(other domain.Entity) bool {
	// TODO implement me
	panic("implement me")
}

func (f *fakeAggregate) ID() string {
	return f.id
}

func (f *fakeAggregate) Version() uint64 {
	return f.version
}

func (f fakeAggregate) withID(id string) fakeAggregate {
	return fakeAggregate{
		id:      id,
		version: f.version,
	}
}

func (f fakeAggregate) withVersion(version uint64) fakeAggregate {
	return fakeAggregate{
		id:      f.id,
		version: version,
	}
}

func (f fakeAggregate) withEvents(events ...domain.Event) fakeAggregate {
	aggregate := fakeAggregate{
		id:      f.id,
		version: f.version,
	}
	for _, event := range events {
		aggregate.Apply(event)
	}
	return aggregate
}
