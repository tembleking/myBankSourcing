package account_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

var _ = Describe("Account Service", func() {
	var (
		accountService    *account.AccountService
		accountRepository *fakeAccountRepository
	)

	BeforeEach(func(ctx context.Context) {
		accountRepository = &fakeAccountRepository{accounts: map[string]*account.Account{}}
		eventStore := persistence.NewEventStoreBuilder(sqlite.InMemory()).Build()
		accountService = account.NewAccountService(eventStore, accountRepository)
	})

	It("opens the account", func(ctx context.Context) {
		accountCreated, err := accountService.OpenAccount(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(accountCreated).ToNot(BeNil())
		Expect(accountCreated.Balance()).To(Equal(0))
		Expect(accountCreated.IsOpen()).To(BeTrue())
		Expect(accountRepository.GetByID(ctx, accountCreated.ID())).To(Equal(accountCreated))
	})

	It("adds money to the account", func(ctx context.Context) {
		accountCreated, err := accountService.OpenAccount(ctx)
		Expect(err).ToNot(HaveOccurred())

		amount := 100
		accountUpdated, err := accountService.DepositMoneyIntoAccount(ctx, accountCreated.ID(), amount)
		Expect(err).ToNot(HaveOccurred())

		Expect(accountUpdated.Balance()).To(Equal(amount))
	})

	It("withdraws money from the account", func(ctx context.Context) {
		accountCreated, err := accountService.OpenAccount(ctx)
		Expect(err).ToNot(HaveOccurred())

		amount := 100
		_, err = accountService.DepositMoneyIntoAccount(ctx, accountCreated.ID(), amount)
		Expect(err).ToNot(HaveOccurred())

		amount = 25
		accountUpdated, err := accountService.WithdrawMoneyFromAccount(ctx, accountCreated.ID(), amount)
		Expect(err).ToNot(HaveOccurred())

		Expect(accountUpdated.Balance()).To(Equal(75))
	})

	It("closes the account", func(ctx context.Context) {
		accountCreated, err := accountService.OpenAccount(ctx)
		Expect(err).ToNot(HaveOccurred())

		updatedAccount, err := accountService.CloseAccount(ctx, accountCreated.ID())
		Expect(err).ToNot(HaveOccurred())

		Expect(updatedAccount.IsOpen()).To(BeFalse())
	})

	It("transfers money between accounts", func(ctx context.Context) {
		firstAccount, err := accountService.OpenAccount(ctx)
		Expect(err).ToNot(HaveOccurred())

		secondAccount, err := accountService.OpenAccount(ctx)
		Expect(err).ToNot(HaveOccurred())

		amountToAdd := 100
		_, err = accountService.DepositMoneyIntoAccount(ctx, firstAccount.ID(), amountToAdd)
		Expect(err).ToNot(HaveOccurred())

		amountToTransfer := 25
		modifiedFirstAccount, err := accountService.TransferMoney(ctx, firstAccount.ID(), secondAccount.ID(), amountToTransfer)
		Expect(err).ToNot(HaveOccurred())
		Expect(modifiedFirstAccount.Balance()).To(Equal(75))
	})
})

type fakeAccountRepository struct {
	accounts map[string]*account.Account
}

func (f *fakeAccountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	if account, ok := f.accounts[id]; ok {
		return account, nil
	}
	return nil, fmt.Errorf("not found")
}

func (f *fakeAccountRepository) Save(ctx context.Context, aggregate *account.Account) error {
	f.accounts[aggregate.ID()] = aggregate
	return nil
}
