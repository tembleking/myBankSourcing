package services_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence/inmemory/account"
)

var _ = Describe("OpenAccount", func() {
	var (
		accountService     *services.OpenAccountService
		inMemoryRepository *account.Repository
	)

	BeforeEach(func() {
		inMemoryRepository = account.NewRepository()
		accountService = services.NewCreateAccountService(inMemoryRepository)
	})

	It("opens the account", func() {
		accountCreated, err := accountService.OpenAccount(context.Background())

		Expect(err).ToNot(HaveOccurred())
		Expect(accountCreated).ToNot(BeNil())
		Expect(accountCreated.Balance()).To(Equal(0))
		Expect(accountCreated.IsOpen()).To(BeTrue())
		Expect(inMemoryRepository.GetAccount(context.Background(), accountCreated.ID())).To(Equal(accountCreated))
	})
})
