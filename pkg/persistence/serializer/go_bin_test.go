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
		serialize, err := ser.Serialize(someEvents())
		Expect(err).ToNot(HaveOccurred())
		Expect(serialize).ToNot(BeEmpty())
		Expect(serialize).To(BeAssignableToTypeOf([]byte{}))

		deserialize, err := ser.Deserialize(serialize)
		Expect(err).ToNot(HaveOccurred())

		Expect(deserialize).To(Equal(someEvents()))
	})

	It("serializes consistently", func() {
		first, err := ser.Serialize(someEvents())
		Expect(err).ToNot(HaveOccurred())

		second, err := ser.Serialize(someEvents())
		Expect(err).ToNot(HaveOccurred())

		Expect(first).To(Equal(second))
	})
})

func someEvents() []domain.Event {
	return []domain.Event{
		&account.AmountAdded{
			Quantity: 50,
			Balance:  134,
		},
		&account.AmountWithdrawn{
			Quantity: 35,
			Balance:  99,
		},
		&account.TransferenceSent{
			From:     "some-id",
			To:       "other-id",
			Quantity: 10,
			Balance:  89,
		},
		&account.TransferenceReceived{
			From:     "some-id",
			To:       "other-id",
			Quantity: 10,
			Balance:  99,
		},
	}
}
