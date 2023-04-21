package outbox_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/outbox"
	outboxmocks "github.com/tembleking/myBankSourcing/pkg/outbox/mocks"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/mocks"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

var _ = Describe("TransactionalOutbox", func() {
	var (
		ctx                    context.Context
		ctrl                   *gomock.Controller
		appendOnlyStore        *mocks.MockAppendOnlyStore
		messageBrokerPublisher *outboxmocks.MockPublishableMessageBroker
		transactionalOutbox    *outbox.TransactionalOutbox
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		appendOnlyStore = mocks.NewMockAppendOnlyStore(ctrl)
		messageBrokerPublisher = outboxmocks.NewMockPublishableMessageBroker(ctrl)

		transactionalOutbox = outbox.NewTransactionalOutboxBuilder(messageBrokerPublisher).
			WithAppendOnlyStore(appendOnlyStore).
			Build()
	})

	When("there are undispatched events", func() {
		It("dispatches all the undispatched events", func() {
			appendOnlyStore.EXPECT().
				ReadUndispatchedRecords(ctx).
				Return([]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventName: "AmountAdded", EventData: dataRecordInStore()}}, nil)

			messageBrokerPublisher.EXPECT().
				Publish(ctx, gomock.Any()).
				Return(nil)

			appendOnlyStore.EXPECT().
				MarkRecordsAsDispatched(ctx, []persistence.StreamID{{StreamName: "aggregate-0", StreamVersion: 1}}).
				Return(nil)

			err := transactionalOutbox.DispatchUndispatchedEvents(ctx)

			Expect(err).To(BeNil())
		})

		When("the append only store does not return any undispatched events", func() {
			It("does nothing", func() {
				appendOnlyStore.EXPECT().
					ReadUndispatchedRecords(ctx).
					Return(nil, nil)

				err := transactionalOutbox.DispatchUndispatchedEvents(ctx)

				Expect(err).To(BeNil())
			})
		})

		When("the publish to the message broker fails", func() {
			It("does not mark the events as dispatched", func() {
				appendOnlyStore.EXPECT().
					ReadUndispatchedRecords(ctx).
					Return([]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventName: "AmountAdded", EventData: dataRecordInStore()}}, nil)

				messageBrokerPublisher.EXPECT().
					Publish(ctx, gomock.Any()).
					Return(errors.New("some error"))

				err := transactionalOutbox.DispatchUndispatchedEvents(ctx)

				Expect(err).To(MatchError(errors.New("some error")))
			})
		})

		When("saving the dispatched events fails", func() {
			It("returns the error", func() {
				appendOnlyStore.EXPECT().
					ReadUndispatchedRecords(ctx).
					Return([]persistence.StoredStreamEvent{{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventName: "AmountAdded", EventData: dataRecordInStore()}}, nil)

				messageBrokerPublisher.EXPECT().
					Publish(ctx, gomock.Any()).
					Return(nil)

				appendOnlyStore.EXPECT().
					MarkRecordsAsDispatched(ctx, []persistence.StreamID{{StreamName: "aggregate-0", StreamVersion: 1}}).
					Return(errors.New("some error"))

				err := transactionalOutbox.DispatchUndispatchedEvents(ctx)

				Expect(err).To(MatchError(errors.New("some error")))
			})
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
