package transfer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/tembleking/myBankSourcing/pkg/transfer"
	. "github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("Transfer", func() {
	It("is equal to itself", func() {
		transfer := RequestTransfer("fromAccount", "toAccount", 100)

		Expect(transfer).To(BeAnEntityEqualTo(transfer))
	})

	It("is created correctly", func() {
		transfer := RequestTransfer("fromAccount", "toAccount", 100)

		Expect(transfer).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"TransferID":  Not(BeEmpty()),
			"FromAccount": Equal("fromAccount"),
			"ToAccount":   Equal("toAccount"),
			"Amount":      Equal(100),
		})))
	})
})
