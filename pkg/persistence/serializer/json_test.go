package serializer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

var _ = Describe("JSON", func() {
	var ser *serializer.JSON
	BeforeEach(func() {
		ser = &serializer.JSON{}
	})

	It("serializes and deserializes the event", func() {
		serialize, err := ser.SerializeDomainEvent(anEvent())
		Expect(err).ToNot(HaveOccurred())
		Expect(serialize).ToNot(BeEmpty())
		Expect(serialize).To(BeAssignableToTypeOf([]byte{}))

		deserialize, err := ser.DeserializeDomainEvent("AmountDeposited", serialize)
		Expect(err).ToNot(HaveOccurred())
		Expect(deserialize).To(Equal(anEvent()))
	})

	It("returns the content type", func() {
		Expect(ser.ContentType()).To(Equal("application/json"))
	})
})
