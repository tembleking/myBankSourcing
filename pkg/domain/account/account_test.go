package account_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

var _ = Describe("Account", func() {
	var (
		acc *account.Account
	)

	BeforeEach(func() {
		acc = account.NewAccount("some-id")
	})

	When("adding money to the account", func() {
		It("should add the money successfully", func() {
			err := acc.AddMoney(50)

			Expect(err).ToNot(HaveOccurred())
			Expect(acc.Balance()).To(Equal(50))
		})

		When("the account already had money", func() {
			It("should return the total balance after adding more money", func() {
				_ = acc.AddMoney(50)

				err := acc.AddMoney(50)

				Expect(err).ToNot(HaveOccurred())
				Expect(acc.Balance()).To(Equal(100))
			})
		})

		When("trying to add a negative amount", func() {
			It("fails", func() {
				err := acc.AddMoney(-1)

				Expect(err).To(MatchError(account.ErrAddMoneyQuantityCannotBeNegative))
			})
		})
	})

	When("removing money from an account", func() {
		When("the account already had money", func() {
			It("subtracts the money", func() {
				_ = acc.AddMoney(50)

				err := acc.WithdrawalMoney(30)

				Expect(err).ToNot(HaveOccurred())
				Expect(acc.Balance()).To(Equal(20))
			})
		})

		When("the account has less money than the amount to withdrawn", func() {
			It("returns an error", func() {
				_ = acc.AddMoney(50)

				err := acc.WithdrawalMoney(51)

				Expect(err).To(MatchError(account.ErrBalanceIsNotEnoughForWithdrawal))
			})
		})
	})

	When("transferring money to another account", func() {
		It("transfers the money correctly", func() {
			origin := account.NewAccount("origin")
			_ = origin.AddMoney(100)
			destination := account.NewAccount("destination")
			_ = destination.AddMoney(30)
			amount := 50

			err := origin.TransferMoney(amount, destination)

			Expect(err).ToNot(HaveOccurred())
			Expect(origin.Balance()).To(Equal(50))
			Expect(destination.Balance()).To(Equal(80))
		})

		When("the origin account has less money than the amount to transfer", func() {
			It("returns an error", func() {
				origin := account.NewAccount("origin")
				_ = origin.AddMoney(30)
				destination := account.NewAccount("destination")
				amount := 50

				err := origin.TransferMoney(amount, destination)

				Expect(err).To(MatchError(account.ErrBalanceIsNotEnoughForTransfer))
			})
		})
	})
})
