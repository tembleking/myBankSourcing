package persistence_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/mocks"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
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
		eventStore = persistence.NewEventStoreBuilder(sqlite.InMemory()).WithAppendOnlyStore(appendOnlyStore).WithClock(&stubClock{}).Build()
	})

	It("should be able to load an event stream", func() {
		appendOnlyStore.EXPECT().ReadRecords(ctx, persistence.StreamName("aggregate-0")).Return(
			[]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventData: dataRecordInStore()}},
			nil,
		)

		stream, err := eventStore.LoadEventStream(ctx, "aggregate-0")

		Expect(err).To(BeNil())
		Expect(stream).To(Equal([]persistence.StreamEvent{{
			ID:    persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1},
			Event: &account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
		}}))
	})

	It("should be able to append to an event stream", func() {
		appendOnlyStore.EXPECT().Append(
			ctx,
			persistence.StoredStreamEvent{
				ID:        persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0},
				EventName: "AmountAdded",
				EventData: dataRecordInStore(),
			},
			persistence.StoredStreamEvent{
				ID:        persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0},
				EventName: "AmountAdded",
				EventData: dataRecordInStore(),
			},
		).Return(nil)

		anAggregate := fakeAggregate{}.withID("aggregate-0").withVersion(1).withEvents(
			&account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
		)
		anotherAggregate := fakeAggregate{}.withID("aggregate-1").withVersion(1).withEvents(
			&account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
		)
		err := eventStore.AppendToStream(ctx, &anAggregate, &anotherAggregate)
		Expect(err).To(BeNil())
	})

	When("asking for all events", func() {
		It("returns all events", func() {
			appendOnlyStore.EXPECT().
				ReadAllRecords(ctx).
				Return([]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventData: dataRecordInStore()}}, nil)

			stream, err := eventStore.LoadAllEvents(ctx)
			Expect(err).To(BeNil())
			Expect(stream).To(Equal([]persistence.StreamEvent{{
				ID:    persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1},
				Event: &account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
			}}))
		})
	})
})

func dataRecordInStore() []byte {
	serializer := &serializer.GoBinarySerializer{}
	data, err := serializer.SerializeDomainEvent(
		&account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
	)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return data
}

type stubClock struct {
}

func (f *stubClock) Now() time.Time {
	return time.Time{}
}

type fakeAggregate struct {
	id      string
	version uint64
	events  []domain.Event
}

// Events implements domain.Aggregate.
func (f *fakeAggregate) Events() []domain.Event {
	return f.events
}

// ID implements domain.Aggregate.
func (f *fakeAggregate) ID() string {
	return f.id
}

// Version implements domain.Aggregate.
func (f *fakeAggregate) Version() uint64 {
	return f.version
}

func (f *fakeAggregate) ClearEvents() {
	f.events = nil
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
	return fakeAggregate{
		id:      f.id,
		version: f.version,
		events:  events,
	}
}
