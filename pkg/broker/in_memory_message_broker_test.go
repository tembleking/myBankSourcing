package broker_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/broker"
	"github.com/tembleking/myBankSourcing/pkg/outbox"
)

var _ = Describe("InMemoryMessageBroker", func() {
	It("publishes the message", func() {
		inMemoryBroker := broker.NewInMemoryMessageBroker()

		err := inMemoryBroker.Publish(context.Background(), []byte("message"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("should receive the message", func() {
		inMemoryBroker := broker.NewInMemoryMessageBroker()
		go inMemoryBroker.StartServing(context.Background())

		received := false
		err := inMemoryBroker.Subscribe(context.Background(), func(message outbox.AcknowledgableMessage) {
			Expect(message.Data()).To(Equal([]byte("message")))
			Expect(message.Ack()).ToNot(HaveOccurred())
			received = true
		})
		Expect(err).ToNot(HaveOccurred())

		err = inMemoryBroker.Publish(context.Background(), []byte("message"))
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool { return received }).Should(BeTrue())
	})

	When("it's not acked", func() {
		It("should receive the message again", func() {
			inMemoryBroker := broker.NewInMemoryMessageBroker()
			go inMemoryBroker.StartServing(context.Background())

			receivedTimes := 0
			err := inMemoryBroker.Subscribe(context.Background(), func(message outbox.AcknowledgableMessage) {
				receivedTimes++
				Expect(message.Data()).To(Equal([]byte("message")))
			})
			Expect(err).ToNot(HaveOccurred())

			err = inMemoryBroker.Publish(context.Background(), []byte("message"))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() int { return receivedTimes }).Should(BeNumerically(">=", 2))
		})
	})

	When("it's acked", func() {
		It("should not receive the message again", func() {
			inMemoryBroker := broker.NewInMemoryMessageBroker()
			go inMemoryBroker.StartServing(context.Background())

			receivedTimes := 0
			err := inMemoryBroker.Subscribe(context.Background(), func(message outbox.AcknowledgableMessage) {
				receivedTimes++
				Expect(message.Data()).To(Equal([]byte("message")))
				Expect(message.Ack()).ToNot(HaveOccurred())
			})
			Expect(err).ToNot(HaveOccurred())

			err = inMemoryBroker.Publish(context.Background(), []byte("message"))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() int { return receivedTimes }).Should(Equal(1))
			Consistently(func() int { return receivedTimes }).Should(Equal(1))
		})
	})

	When("there are multiple receivers, and only one of them acks", func() {
		It("should receive the message again only for the non-acked receivers", func() {
			inMemoryBroker := broker.NewInMemoryMessageBroker()
			go inMemoryBroker.StartServing(context.Background())

			receivedTimesAcker := 0
			err := inMemoryBroker.Subscribe(context.Background(), func(message outbox.AcknowledgableMessage) {
				receivedTimesAcker++
				Expect(message.Data()).To(Equal([]byte("message")))
				Expect(message.Ack()).ToNot(HaveOccurred())
			})
			Expect(err).ToNot(HaveOccurred())

			receivedTimesNonAcker := 0
			err = inMemoryBroker.Subscribe(context.Background(), func(message outbox.AcknowledgableMessage) {
				receivedTimesNonAcker++
				Expect(message.Data()).To(Equal([]byte("message")))
			})
			Expect(err).ToNot(HaveOccurred())

			err = inMemoryBroker.Publish(context.Background(), []byte("message"))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() int { return receivedTimesAcker }).Should(Equal(1))
			Eventually(func() int { return receivedTimesNonAcker }).Should(BeNumerically(">=", 2))
			Consistently(func() int { return receivedTimesAcker }).Should(Equal(1))
			Consistently(func() int { return receivedTimesNonAcker }).Should(BeNumerically(">=", 2))
		})
	})
})
