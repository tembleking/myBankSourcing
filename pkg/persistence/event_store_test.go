package persistence_test

import (
	"context"

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
		eventStore = persistence.NewEventStore(binarySerializer, binarySerializer, appendOnlyStore)
	})

	It("should be able to load an event stream", func() {
		appendOnlyStore.EXPECT().ReadRecords(ctx, "aggregate-0", uint64(0), uint64(0)).Return(
			[]persistence.DataWithVersion{{Version: 1, Data: dataRecordInStore()}},
			nil,
		)

		stream, err := eventStore.LoadEventStream(ctx, "aggregate-0")

		Expect(err).To(BeNil())
		Expect(stream).To(Equal(&persistence.EventStream{
			Name:    "aggregate-0",
			Version: 1,
			Events:  []domain.Event{&account.AmountAdded{BaseEvent: domain.BaseEvent{EventID: "event-0"}, Quantity: 10, Balance: 10}},
		}))
	})

	It("should be able to load an event stream subset", func() {
		appendOnlyStore.EXPECT().ReadRecords(ctx, "aggregate-0", uint64(1), uint64(2)).Return(
			[]persistence.DataWithVersion{{Version: 2, Data: dataRecordInStore()}},
			nil,
		)

		stream, err := eventStore.LoadEventStreamSubset(ctx, "aggregate-0", 1, 2)

		Expect(err).To(BeNil())
		Expect(stream).To(Equal(&persistence.EventStream{
			Name:    "aggregate-0",
			Version: 2,
			Events:  []domain.Event{&account.AmountAdded{BaseEvent: domain.BaseEvent{EventID: "event-0"}, Quantity: 10, Balance: 10}},
		}))
	})

	It("should be able to append to an event stream", func() {
		appendOnlyStore.EXPECT().Append(ctx, "aggregate-0", dataRecordInStore(), uint64(1)).Return(nil)

		err := eventStore.AppendToStream(ctx, "aggregate-0", 1, []domain.Event{
			&account.AmountAdded{
				BaseEvent: domain.BaseEvent{EventID: "event-0"},
				Quantity:  10,
				Balance:   10,
			},
		})

		Expect(err).To(BeNil())
	})

	It("loads the events from all the streams", func() {
		appendOnlyStore.EXPECT().ReadAllRecords(ctx, uint64(0), uint64(0)).Return([]persistence.DataWithNameAndVersion{
			{Name: "aggregate-0", Version: 4, Data: dataRecordInStore()},
			{Name: "aggregate-1", Version: 7, Data: dataRecordInStore()},
		}, nil)

		streams, err := eventStore.LoadAllEventStreams(ctx)

		Expect(err).To(BeNil())
		Expect(streams).To(HaveLen(2))
		Expect(streams[0]).To(Equal(&persistence.EventStream{
			Name:    "aggregate-0",
			Version: 4,
			Events:  []domain.Event{&account.AmountAdded{BaseEvent: domain.BaseEvent{EventID: "event-0"}, Quantity: 10, Balance: 10}},
		}))
		Expect(streams[1]).To(Equal(&persistence.EventStream{
			Name:    "aggregate-1",
			Version: 7,
			Events:  []domain.Event{&account.AmountAdded{BaseEvent: domain.BaseEvent{EventID: "event-0"}, Quantity: 10, Balance: 10}},
		}))
	})
})

func dataRecordInStore() []byte {
	serializer := &serializer.GoBinarySerializer{}
	data, err := serializer.Serialize([]domain.Event{
		&account.AmountAdded{
			BaseEvent: domain.BaseEvent{EventID: "event-0"},
			Quantity:  10,
			Balance:   10,
		},
	},
	)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return data
}
