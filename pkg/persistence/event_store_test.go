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
)

var _ = Describe("EventStore", func() {
	var (
		ctx              context.Context
		binarySerializer *serializer.GoBinarySerializer
		ctrl             *gomock.Controller
		appendOnlyStore  *mocks.MockAppendOnlyStore
		eventStore       *persistence.EventStore
	)

	BeforeEach(func() {
		ctx = context.Background()
		binarySerializer = &serializer.GoBinarySerializer{}
		ctrl = gomock.NewController(GinkgoT())
		appendOnlyStore = mocks.NewMockAppendOnlyStore(ctrl)
		eventStore = persistence.NewEventStore(binarySerializer, binarySerializer, appendOnlyStore, &stubClock{})
	})

	It("should be able to load an event stream", func() {
		appendOnlyStore.EXPECT().ReadRecords(ctx, "aggregate-0").Return(
			[]persistence.StoredStreamEvent{{StreamVersion: 1, EventData: dataRecordInStore()}},
			nil,
		)

		stream, err := eventStore.LoadEventStream(ctx, "aggregate-0")

		Expect(err).To(BeNil())
		Expect(stream).To(Equal([]persistence.StreamEvent{{
			StreamID:      "aggregate-0",
			StreamVersion: 1,
			Event:         &account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
		}}))
	})

	It("should be able to append to an event stream", func() {
		appendOnlyStore.EXPECT().Append(ctx, persistence.StoredStreamEvent{
			StreamID:      "aggregate-0",
			StreamVersion: 0,
			EventName:     "AmountAdded",
			EventData:     dataRecordInStore(),
		}).Return(nil)

		err := eventStore.AppendToStream(ctx, "aggregate-0", 1, []domain.Event{
			&account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
		})
		Expect(err).To(BeNil())
	})

	When("asking for events by name", func() {
		It("should be able to load events by name", func() {
			appendOnlyStore.EXPECT().
				ReadEventsByName(ctx, "AmountAdded").
				Return([]persistence.StoredStreamEvent{{StreamID: "aggregate-0", StreamVersion: 1, EventData: dataRecordInStore()}}, nil)

			stream, err := eventStore.LoadEventsByName(ctx, "AmountAdded")
			Expect(err).To(BeNil())
			Expect(stream).To(Equal([]persistence.StreamEvent{{
				StreamID:      "aggregate-0",
				StreamVersion: 1,
				Event:         &account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
			}}))
		})
	})

	When("configured with a dispatcher", func() {
		It("sends the events to the dispatcher", func() {
			mockDispatcher := mocks.NewMockEventDispatcher(ctrl)
			mockDispatcher.EXPECT().Dispatch(ctx, []persistence.StreamEvent{{
				StreamID: "aggregate-0", Event: &account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
			}})
			appendOnlyStore.EXPECT().Append(ctx, persistence.StoredStreamEvent{
				StreamID:      "aggregate-0",
				StreamVersion: 0,
				EventName:     "AmountAdded",
				EventData:     dataRecordInStore(),
				HappenedOn:    time.Time{},
			}).Return(nil)

			eventStore.AddDispatchers(mockDispatcher)

			err := eventStore.AppendToStream(ctx, "aggregate-0", 1, []domain.Event{
				&account.AmountAdded{AccountID: "some-account", Quantity: 10, Balance: 10},
			})

			Expect(err).To(BeNil())
		})
	})
})

func dataRecordInStore() []byte {
	serializer := &serializer.GoBinarySerializer{}
	data, err := serializer.Serialize(
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
