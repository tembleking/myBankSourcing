package serializer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

var _ = Describe("GoBinary", func() {
	var ser *serializer.GoBinarySerializer
	BeforeEach(func() {
		ser = &serializer.GoBinarySerializer{}
	})

	It("serializes and deserializes the event", func() {
		serialize, err := ser.SerializeDomainEvent(anEvent())
		Expect(err).ToNot(HaveOccurred())
		Expect(serialize).ToNot(BeEmpty())
		Expect(serialize).To(BeAssignableToTypeOf([]byte{}))

		deserialize, err := ser.DeserializeDomainEvent(serialize)
		Expect(err).ToNot(HaveOccurred())

		Expect(deserialize).To(Equal(anEvent()))
	})

	It("serializes consistently", func() {
		first, err := ser.SerializeDomainEvent(anEvent())
		Expect(err).ToNot(HaveOccurred())

		second, err := ser.SerializeDomainEvent(anEvent())
		Expect(err).ToNot(HaveOccurred())

		Expect(first).To(Equal(second))
	})
})

func anEvent() domain.Event {
	return &account.AmountAdded{
		Quantity: 50,
		Balance:  134,
	}
}
