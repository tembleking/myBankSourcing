package account_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
	"github.com/tembleking/myBankSourcing/test/mother"
)

var _ = Describe("Accounts", func() {
	var eventStore *persistence.ReadOnlyEventStore

	BeforeEach(func() {
		store := persistence.NewEventStoreBuilder(sqlite.InMemory()).Build()

		Expect(store.AppendToStream(context.Background(), mother.AccountOpenWithMovements())).To(Succeed())

		eventStore = store.ReadOnlyEventStore
	})

	When("there are multiple AccountOpened events saved", func() {
		It("returns the accounts found", func() {
			accountsProjection, err := account.NewAccountProjection(eventStore)
			Expect(err).ToNot(HaveOccurred())

			accounts := accountsProjection.Accounts()

			Expect(accounts).To(HaveLen(1))
			Expect(accounts[0].ID()).To(Equal("some-account"))
			Expect(accounts[0].Balance()).To(Equal(5))
			Expect(accounts[0].Version()).To(Equal(uint64(4)))
		})
	})
})
