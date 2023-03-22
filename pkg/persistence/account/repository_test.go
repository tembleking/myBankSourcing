package account_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	account2 "github.com/tembleking/myBankSourcing/pkg/persistence/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/inmemory"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	. "github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("In Memory Repository", func() {
	var repository *account2.Repository

	BeforeEach(func() {
		eventSerializer := &serializer.GoBinarySerializer{}
		eventStore := persistence.NewEventStore(eventSerializer, eventSerializer, inmemory.NewAppendOnlyStore())
		repository = account2.NewRepository(eventStore)
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

	It("lists the accounts", func() {
		oneAccount := someAccountWithMovementsAndID("some-id")
		anotherAccount := someAccountWithMovementsAndID("some-other-id")

		err := repository.SaveAccount(context.Background(), oneAccount)
		Expect(err).ToNot(HaveOccurred())

		err = repository.SaveAccount(context.Background(), anotherAccount)
		Expect(err).ToNot(HaveOccurred())

		accounts, err := repository.ListAccounts(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).To(ConsistOf(BeAnAccountEqualsTo(oneAccount), BeAnAccountEqualsTo(anotherAccount)))
	})
})

func someAccountWithMovements() *account.Account {
	return someAccountWithMovementsAndID("some-id")
}

func someAccountWithMovementsAndID(id string) *account.Account {
	anAccount := account.NewAccount()

	_ = anAccount.OpenAccount(account.ID(id))
	_ = anAccount.AddMoney(50)
	_ = anAccount.WithdrawalMoney(35)
	_ = anAccount.AddMoney(10)

	return anAccount
}
