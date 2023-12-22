package transfer_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
	"github.com/tembleking/myBankSourcing/pkg/transfer"
	"github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("Repository", func() {
	var repository domain.Repository[*transfer.Transfer]

	BeforeEach(func() {
		repository = transfer.NewRepository(persistence.NewEventStoreBuilder(sqlite.InMemory()).Build())
	})

	It("saves an transfer and retrieves it", func(ctx context.Context) {
		acc := transfer.RequestTransfer("from-account", "to-account", 50)

		Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
		Expect(repository.GetByID(ctx, acc.ID())).To(matchers.BeAnEntityEqualTo(acc))
		Expect(repository.GetByID(ctx, acc.ID())).To(matchers.BeAggregateWithTheSameVersionAs(acc))
	})

	When("saving the transfer multiple times", func() {
		It("returns an error", func(ctx context.Context) {
			acc := transfer.RequestTransfer("from-account", "to-account", 50)

			Expect(repository.Save(ctx, acc)).ToNot(HaveOccurred())
			Expect(repository.Save(ctx, acc)).To(HaveOccurred())
		})
	})
})
