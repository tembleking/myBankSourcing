package account_test

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
	"github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("Repository", func() {
	var repository domain.Repository[*account.Account]

	BeforeEach(func() {
		repository = account.NewRepository(persistence.NewEventStoreBuilder(sqlite.InMemory()).Build())
	})

	It("saves an account and retrieves it", func(ctx context.Context) {
		id := uuid.NewString()
		acc := account.OpenAccount(id)
		Expect(acc.AddMoney(50)).ToNot(HaveOccurred())
		Expect(acc.WithdrawMoney(20)).ToNot(HaveOccurred())
		Expect(acc.WithdrawMoney(30)).ToNot(HaveOccurred())
		Expect(acc.CloseAccount()).ToNot(HaveOccurred())

		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
		Expect(repository.GetByID(ctx, id)).To(matchers.BeAnAccountEqualsTo(acc))
		Expect(repository.GetByID(ctx, id)).To(matchers.BeAggregateWithTheSameVersionAs(acc))
	})

	It("saves the same account multiple times and is idempotent", func(ctx context.Context) {
		id := uuid.NewString()
		acc := account.OpenAccount(id)
		Expect(acc.AddMoney(50)).ToNot(HaveOccurred())
		Expect(acc.WithdrawMoney(20)).ToNot(HaveOccurred())
		Expect(acc.WithdrawMoney(30)).ToNot(HaveOccurred())
		Expect(acc.CloseAccount()).ToNot(HaveOccurred())

		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
	})
})
