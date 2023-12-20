package eventbus_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/eventbus"
)

var _ = Describe("Eventbus", func() {
	var (
		eventBus   domain.EventBus
		subscriber *fakeSubscriber
	)

	BeforeEach(func(ctx context.Context) {
		eventBus = eventbus.NewInMemory()
		subscriber = &fakeSubscriber{}

		Expect(eventBus.Subscribe(ctx, subscriber)).To(Succeed())
	})

	It("should publish events to subscribers", func(ctx context.Context) {
		err := eventBus.Publish(ctx, &account.AccountOpened{})
		Expect(err).NotTo(HaveOccurred())

		Expect(subscriber.eventsReceived).To(HaveLen(1))
		Expect(subscriber.eventsReceived[0]).To(BeAssignableToTypeOf(&account.AccountOpened{}))
	})

	It("should send the events in order", MustPassRepeatedly(100), func(ctx context.Context) {
		err := eventBus.Publish(
			ctx,
			&account.AccountOpened{},
			&account.AmountDeposited{},
			&account.AmountWithdrawn{},
			&account.AccountClosed{},
		)
		Expect(err).NotTo(HaveOccurred())

		Expect(subscriber.eventsReceived).To(HaveLen(4))
		Expect(subscriber.eventsReceived[0]).To(BeAssignableToTypeOf(&account.AccountOpened{}))
		Expect(subscriber.eventsReceived[1]).To(BeAssignableToTypeOf(&account.AmountDeposited{}))
		Expect(subscriber.eventsReceived[2]).To(BeAssignableToTypeOf(&account.AmountWithdrawn{}))
		Expect(subscriber.eventsReceived[3]).To(BeAssignableToTypeOf(&account.AccountClosed{}))
	})
})

type fakeSubscriber struct {
	eventsReceived []domain.Event
}

func (f *fakeSubscriber) OnEvent(ctx context.Context, event domain.Event) error {
	f.eventsReceived = append(f.eventsReceived, event)
	return nil
}
