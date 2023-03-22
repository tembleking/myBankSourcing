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
		appendOnlyStore.EXPECT().ReadRecords(ctx, "aggregate-0").Return(
			[]persistence.DataWithVersion{{Version: 1, Data: dataRecordInStore()}},
			nil,
		)

		stream, err := eventStore.LoadEventStream(ctx, "aggregate-0")

		Expect(err).To(BeNil())
		Expect(stream).To(Equal(&persistence.EventStream{
			Name:    "aggregate-0",
			Version: 1,
			Events:  []domain.Event{&account.AmountAdded{Quantity: 10, Balance: 10}},
		}))
	})

	It("should be able to load an event stream subset", func() {
		appendOnlyStore.EXPECT().ReadRecords(ctx, "aggregate-0").Return(
			[]persistence.DataWithVersion{{Version: 2, Data: dataRecordInStore()}},
			nil,
		)

		stream, err := eventStore.LoadEventStreamSubset(ctx, "aggregate-0")

		Expect(err).To(BeNil())
		Expect(stream).To(Equal(&persistence.EventStream{
			Name:    "aggregate-0",
			Version: 2,
			Events:  []domain.Event{&account.AmountAdded{Quantity: 10, Balance: 10}},
		}))
	})

	It("should be able to append to an event stream", func() {
		appendOnlyStore.EXPECT().Append(ctx, "aggregate-0", dataRecordInStore(), uint64(1)).Return(nil)

		err := eventStore.AppendToStream(ctx, "aggregate-0", 1, []domain.Event{
			&account.AmountAdded{Quantity: 10, Balance: 10},
		})
		Expect(err).To(BeNil())
	})
})

func dataRecordInStore() []byte {
	serializer := &serializer.GoBinarySerializer{}
	data, err := serializer.Serialize([]domain.Event{
		&account.AmountAdded{Quantity: 10, Balance: 10},
	})
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return data
}
