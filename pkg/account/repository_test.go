package account_test

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
	"github.com/tembleking/myBankSourcing/test/matchers"
	"github.com/tembleking/myBankSourcing/test/mother"
)

var _ = Describe("Repository", func() {
	var repository domain.Repository[*account.Account]

	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		repository = account.NewRepository(persistence.NewEventStoreBuilder(sqlite.InMemory()).Build())
	})

	It("saves an account and retrieves it", func(ctx context.Context) {
		acc := mother.AccountOpenWithMovements()

		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
		Expect(repository.GetByID(ctx, "some-account")).To(matchers.BeAnEntityEqualTo(acc))
		Expect(repository.GetByID(ctx, "some-account")).To(matchers.BeAggregateWithTheSameVersionAs(acc))
	})

	When("saving the account multiple times", func() {
		It("returns an error", func(ctx context.Context) {
			acc := mother.AccountOpenWithMovements()

			Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
			Expect(repository.Save(ctx, acc)).To(HaveOccurred())
		})
	})

	When("the requested account does not exist", func() {
		It("returns an error", func(ctx context.Context) {
			_, err := repository.GetByID(ctx, "some-account")

			Expect(err).To(MatchError(account.ErrAccountNotFound))
		})
	})
})
