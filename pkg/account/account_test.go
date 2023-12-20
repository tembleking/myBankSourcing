package account_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
)

var _ = Describe("Account", func() {
	var acc *account.Account

	When("the account is not yet open", func() {
		BeforeEach(func() {
			acc = account.NewAccount()
		})

		It("opens the account correctly", func() {
			acc, err := account.OpenAccount("some-id")

			Expect(err).ToNot(HaveOccurred())
			Expect(acc.ID()).To(Equal("some-id"))
			Expect(acc.IsOpen()).To(BeTrue())
		})

		When("opened with an empty ID", func() {
			It("returns an error", func() {
				_, err := account.OpenAccount("")

				Expect(err).To(MatchError("id must not be empty"))
			})
		})

		When("performing any action on the account", func() {
			It("fails if the account is not open", func() {
				Expect(acc.DepositMoney(50)).To(MatchError(account.ErrAccountIsClosed))
			})

			It("fails if the account is not open", func() {
				Expect(acc.WithdrawMoney(50)).To(MatchError(account.ErrAccountIsClosed))
			})
		})
	})

	When("the account is already open", func() {
		BeforeEach(func() {
			var err error
			acc, err = account.OpenAccount("some-id")
			Expect(err).ToNot(HaveOccurred())
		})

		When("depositing money to the account", func() {
			It("should deposit the money successfully", func() {
				err := acc.DepositMoney(50)

				Expect(err).ToNot(HaveOccurred())
				Expect(acc.Balance()).To(Equal(50))
			})

			When("the account already had money", func() {
				It("should return the total balance after adding more money", func() {
					_ = acc.DepositMoney(50)

					err := acc.DepositMoney(50)

					Expect(err).ToNot(HaveOccurred())
					Expect(acc.Balance()).To(Equal(100))
				})
			})

			When("trying to add a negative amount", func() {
				It("fails", func() {
					Expect(acc.DepositMoney(-1)).To(MatchError(account.ErrDepositMoneyQuantityCannotBeNegative))
				})
			})
		})

		When("removing money from an account", func() {
			When("the account already had money", func() {
				It("subtracts the money", func() {
					_ = acc.DepositMoney(50)

					err := acc.WithdrawMoney(30)

					Expect(err).ToNot(HaveOccurred())
					Expect(acc.Balance()).To(Equal(20))
				})
			})

			When("the account has less money than the amount to withdrawn", func() {
				It("returns an error", func() {
					_ = acc.DepositMoney(50)

					err := acc.WithdrawMoney(51)

					Expect(err).To(MatchError(account.ErrBalanceIsNotEnoughForWithdrawal))
				})
			})
		})
	})

	When("still contains balance", func() {
		It("cannot be closed", func() {
			acc, _ := account.OpenAccount("some-id")
			_ = acc.DepositMoney(50)

			err := acc.CloseAccount()

			Expect(err).To(MatchError(account.ErrAccountCannotBeClosedWithBalance))
		})
	})
})
