package transfer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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

		Expect(transfer.ID()).ToNot(BeEmpty())
		Expect(transfer.FromAccount()).To(Equal("fromAccount"))
		Expect(transfer.ToAccount()).To(Equal("toAccount"))
		Expect(transfer.Amount()).To(Equal(100))
	})
})
