package services_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	accountpersistence "github.com/tembleking/myBankSourcing/pkg/persistence/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/inmemory"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	. "github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("Account", func() {
	var (
		accountService    *services.AccountService
		accountRepository *accountpersistence.Repository
	)

	BeforeEach(func() {
		serializer := &serializer.GoBinarySerializer{}
		eventStore := persistence.NewEventStore(serializer, serializer, inmemory.NewAppendOnlyStore())
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

	It("lists the accounts", func() {
		oneAccount, err := accountService.OpenAccount(context.Background())
		Expect(err).ToNot(HaveOccurred())

		anotherAccount, err := accountService.OpenAccount(context.Background())
		Expect(err).ToNot(HaveOccurred())

		accounts, err := accountService.ListAccounts(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).To(ConsistOf(BeAnAccountEqualsTo(oneAccount), BeAnAccountEqualsTo(anotherAccount)))
		Expect(oneAccount).ToNot(BeAnAccountEqualsTo(anotherAccount))
	})
})
