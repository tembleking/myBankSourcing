package account_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

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
			Expect(accounts[0]).To(MatchFields(IgnoreExtras, Fields{
				"AccountID": Equal("some-account"),
				"Balance":   Equal(5),
				"Movements": ConsistOf(
					MatchFields(IgnoreExtras, Fields{
						"Type":             Equal("Deposit"),
						"Amount":           Equal(50),
						"ResultingBalance": Equal(50),
						"Timestamp":        BeTemporally("~", time.Now(), 2*time.Second),
					}),
					MatchFields(IgnoreExtras, Fields{
						"Type":             Equal("Withdrawal"),
						"Amount":           Equal(30),
						"ResultingBalance": Equal(20),
					}),
					MatchFields(IgnoreExtras, Fields{
						"Type":             Equal("Withdrawal"),
						"Amount":           Equal(15),
						"ResultingBalance": Equal(5),
						"Timestamp":        BeTemporally("~", time.Now(), 2*time.Second),
					}),
				),
			}))
		})

		When("updating an account between refreshes", func() {
			var accountsProjection *account.Projection
			BeforeEach(func(ctx context.Context) {
				var err error
				accountsProjection, err = account.NewAccountProjection(context.Background(), eventStore.ReadOnlyEventStore, 100*time.Millisecond)
				Expect(err).ToNot(HaveOccurred())

				accounts := accountsProjection.Accounts()
				Expect(accounts).To(HaveLen(1))

				events, err := eventStore.LoadEventStream(ctx, accounts[0].AccountID)
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
				Expect(accountsBeforeRefresh[0].Balance).To(Equal(5))

				Eventually(func() int {
					accountsAfterRefresh := accountsProjection.Accounts()
					return accountsAfterRefresh[0].Balance
				}).Should(Equal(105))
			})
		})
	})
})
