package account_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	account2 "github.com/tembleking/myBankSourcing/pkg/persistence/inmemory/account"
)

var _ = Describe("In Memory Repository", func() {
	var (
		repository *account2.Repository
	)

	BeforeEach(func() {
		repository = account2.NewRepository()
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

			Expect(err).To(MatchError("not found"))
		})
	})
})

func someAccountWithMovements() *account.Account {
	anAccount := account.NewAccount()

	_ = anAccount.OpenAccount("some-id")
	_ = anAccount.AddMoney(50)
	_ = anAccount.WithdrawalMoney(35)
	_ = anAccount.AddMoney(10)

	return anAccount
}
