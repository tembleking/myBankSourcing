package transfer_test

import (
	"context"

	. "github.com/onsi/gomega/gstruct"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/transfer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

var _ = Describe("Transfers", func() {
	var eventStore *persistence.EventStore

	BeforeEach(func() {
		eventStore = persistence.NewEventStoreBuilder(sqlite.InMemory()).Build()
		err := eventStore.AppendToStream(context.Background(), someAggregate())
		Expect(err).ToNot(HaveOccurred())
	})

	When("there are multiple TransferSent events saved", func() {
		It("returns the transfers found", func() {
			transfersProjection, err := transfer.NewTransfersProjectionFrom(eventStore)
			Expect(err).ToNot(HaveOccurred())

			transfers := transfersProjection.Transfers()

			Expect(transfers).To(ConsistOf(MatchFields(IgnoreExtras, Fields{
				"From":     Equal("some-account"),
				"To":       Equal("another-account"),
				"Quantity": Equal(50),
			})))
		})
	})
})

func someAggregate() domain.Aggregate {
	aggregate, err := account.OpenAccount("some-account")
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	anotherAggregate, err := account.OpenAccount("another-account")
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	ExpectWithOffset(1, aggregate.DepositMoney(50)).To(Succeed())
	ExpectWithOffset(1, aggregate.DepositMoney(100)).To(Succeed())
	ExpectWithOffset(1, aggregate.TransferMoney(50, anotherAggregate)).To(Succeed())
	ExpectWithOffset(1, aggregate.WithdrawMoney(50)).To(Succeed())

	return aggregate
}
