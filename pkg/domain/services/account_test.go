package services_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	accountpersistence "github.com/tembleking/myBankSourcing/pkg/persistence/account"
)

var _ = Describe("Account", func() {
	var (
		accountService    *services.AccountService
		accountRepository *accountpersistence.Repository
	)

	BeforeEach(func() {
		eventStore := persistence.NewEventStoreBuilder().Build()
		accountRepository = accountpersistence.NewRepository(eventStore)
		accountService = services.NewAccountService(accountRepository)
	})

	It("opens the account", func() {
		accountCreated, err := accountService.OpenAccount(context.Background())

		Expect(err).ToNot(HaveOccurred())
		Expect(accountCreated).ToNot(BeNil())
		Expect(accountCreated.Balance()).To(Equal(0))
		Expect(accountCreated.IsOpen()).To(BeTrue())

		accountSaved, err := accountRepository.GetAccount(context.Background(), accountCreated.ID())
		Expect(err).ToNot(HaveOccurred())
		Expect(accountSaved.ID()).To(Equal(accountCreated.ID()))
		Expect(accountSaved.Balance()).To(Equal(accountCreated.Balance()))
		Expect(accountSaved.IsOpen()).To(Equal(accountCreated.IsOpen()))
	})
})
