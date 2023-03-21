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

	When("the account is not yet open", func() {
		BeforeEach(func() {
			acc = account.NewAccount()
		})

		It("opens the account correctly", func() {
			acc := account.NewAccount()
			err := acc.OpenAccount("some-id")

			Expect(err).ToNot(HaveOccurred())
			Expect(acc.ID()).To(Equal(account.ID("some-id")))
			Expect(acc.IsOpen()).To(BeTrue())
		})

		When("performing any action on the account", func() {
			It("fails if the account is not open", func() {
				err := acc.AddMoney(50)

				Expect(err).To(MatchError(account.ErrAccountIsClosed))
			})

			It("fails if the account is not open", func() {
				err := acc.WithdrawalMoney(50)

				Expect(err).To(MatchError(account.ErrAccountIsClosed))
			})
		})
	})

	When("the account is already open", func() {
		BeforeEach(func() {
			acc = account.NewAccount(&account.AccountOpened{AccountID: "some-id"})
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
	})

	When("transferring money to another account", func() {
		var (
			origin      *account.Account
			destination *account.Account
		)
		BeforeEach(func() {
			origin = account.NewAccount(
				&account.AccountOpened{AccountID: "origin"},
				&account.AmountAdded{Balance: 100, Quantity: 100},
			)
			destination = account.NewAccount(
				&account.AccountOpened{AccountID: "destination"},
				&account.AmountAdded{Balance: 30, Quantity: 30},
			)
		})

		It("transfers the money correctly", func() {
			amount := 50

			err := origin.TransferMoney(amount, destination)

			Expect(err).ToNot(HaveOccurred())
			Expect(origin.Balance()).To(Equal(50))
			Expect(destination.Balance()).To(Equal(80))
		})

		When("the origin account has less money than the amount to transfer", func() {
			It("returns an error", func() {
				amount := 101

				err := origin.TransferMoney(amount, destination)

				Expect(err).To(MatchError(account.ErrBalanceIsNotEnoughForTransfer))
			})
		})
	})
})
