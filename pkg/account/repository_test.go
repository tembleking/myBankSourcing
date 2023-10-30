package account_test

import (
	"context"

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
		repository = account.NewRepository(persistence.NewEventStoreBuilder(sqlite.InMemory()).Build())
	})

	It("saves an account and retrieves it", func(ctx context.Context) {
		acc := mother.AccountOpenWithMovements()

		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
		Expect(repository.GetByID(ctx, "some-account")).To(matchers.BeAnAccountEqualsTo(acc))
		Expect(repository.GetByID(ctx, "some-account")).To(matchers.BeAggregateWithTheSameVersionAs(acc))
	})

	It("saves the same account multiple times and is idempotent", func(ctx context.Context) {
		acc := mother.AccountOpenWithMovements()

		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
	})
})
