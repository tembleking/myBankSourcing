package serializer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

var _ = Describe("Msgpack", func() {
	var ser *serializer.Msgpack
	BeforeEach(func() {
		ser = &serializer.Msgpack{}
	})

	It("serializes and deserializes the event", func() {
		serialize, err := ser.Serialize(anEvent())
		Expect(err).ToNot(HaveOccurred())
		Expect(serialize).ToNot(BeEmpty())
		Expect(serialize).To(BeAssignableToTypeOf([]byte{}))

		deserialize, err := ser.Deserialize(serialize)
		Expect(err).ToNot(HaveOccurred())

		Expect(deserialize).To(Equal(anEvent()))
	})
})
