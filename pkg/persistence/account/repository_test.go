package account_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	accountpersistence "github.com/tembleking/myBankSourcing/pkg/persistence/account"
	. "github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("In Memory Repository", func() {
	var repository *accountpersistence.Repository

	BeforeEach(func() {
		eventStore := persistence.NewEventStoreBuilder().Build()
		repository = accountpersistence.NewRepository(eventStore)
	})

	It("saves the account and retrieves it again", func() {
		savedAccount := someAccountWithMovements()

		err := repository.SaveAccount(context.Background(), savedAccount)
		Expect(err).ToNot(HaveOccurred())

		foundAccount, err := repository.GetAccount(context.Background(), "some-id")
		Expect(err).ToNot(HaveOccurred())
		Expect(foundAccount.ID()).To(Equal(account.ID("some-id")))
		Expect(foundAccount.Balance()).To(Equal(25))
	})

	When("the account does not exist", func() {
		It("returns an error", func() {
			_, err := repository.GetAccount(context.Background(), "unknown")

			Expect(err).To(MatchError(ContainSubstring("not found")))
		})
	})

	It("returns all the accounts", func() {
		_ = repository.SaveAccount(context.Background(), someAccountWithMovementsAndID("some-id-1"))
		_ = repository.SaveAccount(context.Background(), someAccountWithMovementsAndID("some-id-2"))
		_ = repository.SaveAccount(context.Background(), someAccountWithMovementsAndID("some-id-3"))

		accounts, err := repository.ListAccounts(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).To(HaveLen(3))
		Expect(accounts).To(ConsistOf(
			BeAnAccountEqualsTo(someAccountWithMovementsAndID("some-id-1")),
			BeAnAccountEqualsTo(someAccountWithMovementsAndID("some-id-2")),
			BeAnAccountEqualsTo(someAccountWithMovementsAndID("some-id-3")),
		))
	})
})

func someAccountWithMovements() *account.Account {
	return someAccountWithMovementsAndID("some-id")
}

func someAccountWithMovementsAndID(id string) *account.Account {
	anAccount := account.OpenAccount(account.ID(id))

	_ = anAccount.AddMoney(50)
	_ = anAccount.WithdrawMoney(35)
	_ = anAccount.AddMoney(10)

	return anAccount
}
