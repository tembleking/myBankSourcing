package account_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	account2 "github.com/tembleking/myBankSourcing/pkg/domain/account"
)

var _ = Describe("Account", func() {
	var (
		acc *account2.Account
	)

	BeforeEach(func() {
		acc = account2.NewAccount("some-id")
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

				Expect(err).To(MatchError(account2.ErrAddMoneyQuantityCannotBeNegative))
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

				Expect(err).To(MatchError(account2.ErrBalanceIsNotEnoughForWithdrawal))
			})
		})
	})
})
