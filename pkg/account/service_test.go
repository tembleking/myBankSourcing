package account_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/inmemory"
	"github.com/tembleking/myBankSourcing/pkg/transfer"
	. "github.com/tembleking/myBankSourcing/test/matchers"
)

var _ = Describe("Account Service", func() {
	var (
		accountService     *account.Service
		accountRepository  *inmemory.Repository[*account.Account]
		transferRepository *inmemory.Repository[*transfer.Transfer]
	)

	BeforeEach(func(ctx context.Context) {
		accountRepository = inmemory.NewRepository[*account.Account]()
		transferRepository = inmemory.NewRepository[*transfer.Transfer]()
		accountService = account.NewAccountService(accountRepository, transferRepository)
	})

	It("opens the account", func(ctx context.Context) {
		accountCreated, err := accountService.OpenAccount(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(accountCreated).ToNot(BeNil())
		Expect(accountCreated.Balance()).To(Equal(0))
		Expect(accountCreated.IsOpen()).To(BeTrue())
		Expect(accountRepository.GetByID(ctx, accountCreated.ID())).To(Equal(accountCreated))
	})

	It("opens the account", func(ctx context.Context) {
		err := accountService.OnCommand(ctx, &account.OpenNewAccount{ID: "some-account-id"})
		Expect(err).ToNot(HaveOccurred())

		accountCreated, err := accountRepository.GetByID(ctx, "some-account-id")
		Expect(err).ToNot(HaveOccurred())
		Expect(accountCreated).ToNot(BeNil())
		Expect(accountCreated.Balance()).To(Equal(0))
		Expect(accountCreated.IsOpen()).To(BeTrue())
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

	When("the account already exists", func() {
		var (
			origin      *account.Account
			destination *account.Account
		)

		BeforeEach(func(ctx context.Context) {
			var err error
			origin, err = accountService.OpenAccount(ctx)
			Expect(err).ToNot(HaveOccurred())

			destination, err = accountService.OpenAccount(ctx)
			Expect(err).ToNot(HaveOccurred())

			amount := 100
			origin, err = accountService.DepositMoneyIntoAccount(ctx, origin.ID(), amount)
			Expect(err).ToNot(HaveOccurred())
		})

		It("creates a transfer request", func(ctx context.Context) {
			amountToTransfer := 50
			transfer, err := accountService.TransferMoney(ctx, origin.ID(), destination.ID(), amountToTransfer)

			Expect(err).ToNot(HaveOccurred())
			Expect(transfer.ID()).ToNot(BeEmpty())
			Expect(transfer.Amount()).To(Equal(50))
			Expect(transfer.FromAccount()).To(Equal(origin.ID()))
			Expect(transfer.ToAccount()).To(Equal(destination.ID()))

			Expect(transferRepository.GetByID(ctx, transfer.ID())).To(BeAnEntityEqualTo(transfer))
		})

		When("the transfer has been created", func() {
			var transferRequested *transfer.Transfer

			BeforeEach(func(ctx context.Context) {
				var err error
				transferRequested, err = accountService.TransferMoney(ctx, origin.ID(), destination.ID(), 50)
				Expect(err).ToNot(HaveOccurred())
			})

			It("sends the transfer", func(ctx context.Context) {
				err := accountService.SendTransfer(ctx, transferRequested.ID())
				Expect(err).ToNot(HaveOccurred())

				originModified, err := accountRepository.GetByID(ctx, origin.ID())
				Expect(err).ToNot(HaveOccurred())
				Expect(originModified.Balance()).To(Equal(50))
			})

			It("receives the transfer", func(ctx context.Context) {
				err := accountService.ReceiveTransfer(ctx, transferRequested.ID())
				Expect(err).ToNot(HaveOccurred())

				destinationModified, err := accountRepository.GetByID(ctx, destination.ID())
				Expect(err).ToNot(HaveOccurred())
				Expect(destinationModified.Balance()).To(Equal(50))
			})

			When("the transfer has already been sent", func() {
				BeforeEach(func(ctx context.Context) {
					Expect(accountService.SendTransfer(ctx, transferRequested.ID())).To(Succeed())
				})

				It("rollsback the transfer", func(ctx context.Context) {
					err := accountService.RollbackTransfer(ctx, transferRequested.ID())
					Expect(err).ToNot(HaveOccurred())

					originModified, err := accountRepository.GetByID(ctx, origin.ID())
					Expect(err).ToNot(HaveOccurred())
					Expect(originModified.Balance()).To(Equal(100))
				})

				It("completes the transfer", func(ctx context.Context) {
					err := accountService.CompleteTransfer(ctx, transferRequested.ID())
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
