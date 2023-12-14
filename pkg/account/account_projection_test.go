package account_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
	"github.com/tembleking/myBankSourcing/test/mother"
)

var _ = Describe("Accounts", func() {
	var eventStore *persistence.EventStore

	BeforeEach(func() {
		eventStore = persistence.NewEventStoreBuilder(sqlite.InMemory()).Build()
		Expect(eventStore.AppendToStream(context.Background(), mother.AccountOpenWithMovements())).To(Succeed())
	})

	When("there are multiple AccountOpened events saved", func() {
		It("returns the accounts found", func(ctx context.Context) {
			accountsProjection, err := account.NewAccountProjection(ctx, eventStore.ReadOnlyEventStore, time.Second)
			Expect(err).ToNot(HaveOccurred())

			accounts := accountsProjection.Accounts()

			Expect(accounts).To(HaveLen(1))
			Expect(accounts[0].ID()).To(Equal("some-account"))
			Expect(accounts[0].Balance()).To(Equal(5))
			Expect(accounts[0].Version()).To(Equal(uint64(4)))
		})

		When("updating an account between refreshes", func() {
			var accountsProjection *account.Projection
			BeforeEach(func(ctx context.Context) {
				var err error
				accountsProjection, err = account.NewAccountProjection(context.Background(), eventStore.ReadOnlyEventStore, 100*time.Millisecond)
				Expect(err).ToNot(HaveOccurred())

				accounts := accountsProjection.Accounts()
				Expect(accounts).To(HaveLen(1))
				Expect(accounts[0].Balance()).To(Equal(5))

				events, err := eventStore.LoadEventStream(ctx, accounts[0].ID())
				Expect(err).ToNot(HaveOccurred())

				accountToUpdate := account.NewAccount()
				accountToUpdate.LoadFromHistory(events...)
				Expect(accountToUpdate.DepositMoney(100)).To(Succeed())

				err = eventStore.AppendToStream(ctx, accountToUpdate)
				Expect(err).ToNot(HaveOccurred())
			})

			It("refreshes the projection after some time", func(ctx context.Context) {
				accountsBeforeRefresh := accountsProjection.Accounts()
				Expect(accountsBeforeRefresh).To(HaveLen(1))
				Expect(accountsBeforeRefresh[0].Balance()).To(Equal(5))

				Eventually(func() int {
					accountsAfterRefresh := accountsProjection.Accounts()
					return accountsAfterRefresh[0].Balance()
				}).Should(Equal(105))
			})
		})
	})

})
