package account_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/transfer"
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
					Expect(acc.DepositMoney(-1)).To(MatchError(account.ErrQuantityCannotBeNegative))
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

					Expect(err).To(MatchError(account.ErrBalanceIsNotEnough))
				})
			})

			When("the removal of the money is negative", func() {
				It("fails", func() {
					Expect(acc.WithdrawMoney(-1)).To(MatchError(account.ErrQuantityCannotBeNegative))
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

	When("transferring money to another account", func() {
		var (
			origin      *account.Account
			destination *account.Account
		)
		BeforeEach(func() {
			var err error
			origin, err = account.OpenAccount("origin")
			Expect(err).ToNot(HaveOccurred())
			Expect(origin.DepositMoney(100)).To(Succeed())

			destination, err = account.OpenAccount("destination")
			Expect(err).ToNot(HaveOccurred())
			Expect(destination.DepositMoney(30)).To(Succeed())
		})

		It("creates the transfer from one account to another", func() {
			amount := 50

			transfer, err := origin.TransferMoney(amount, destination)

			Expect(err).ToNot(HaveOccurred())
			Expect(transfer.ID()).ToNot(BeEmpty())
			Expect(transfer.FromAccount()).To(Equal(origin.ID()))
			Expect(transfer.ToAccount()).To(Equal(destination.ID()))
			Expect(transfer.Amount()).To(Equal(amount))
		})

		When("the account is the same", func() {
			It("fails", func() {
				_, err := origin.TransferMoney(50, origin)

				Expect(err).To(MatchError(account.ErrCannotTransferToSameAccount))
			})
		})

		When("the origin account is closed", func() {
			BeforeEach(func() {
				Expect(origin.WithdrawMoney(origin.Balance())).To(Succeed())
				Expect(origin.CloseAccount()).To(Succeed())
			})

			It("cannot transfer any money", func() {
				amount := 50

				_, err := origin.TransferMoney(amount, destination)
				Expect(err).To(MatchError(account.ErrAccountIsClosed))
			})
		})

		When("the destination account is closed", func() {
			BeforeEach(func() {
				Expect(destination.WithdrawMoney(destination.Balance())).To(Succeed())
				Expect(destination.CloseAccount()).To(Succeed())
			})

			It("cannot transfer any money", func() {
				amount := 50

				_, err := origin.TransferMoney(amount, destination)
				Expect(err).To(MatchError(account.ErrAccountIsClosed))
			})
		})

		When("the origin account doesn't have enough balance", func() {
			It("fails to transfer the money", func() {
				tooMuchAmount := 200

				_, err := origin.TransferMoney(tooMuchAmount, destination)
				Expect(err).To(MatchError(account.ErrBalanceIsNotEnough))
			})
		})

		When("trying to transfer negative amount", func() {
			It("fails", func() {
				_, err := origin.TransferMoney(-1, destination)
				Expect(err).To(MatchError(account.ErrQuantityCannotBeNegative))
			})
		})
	})

	When("assigning a transfer to an account", func() {
		var (
			origin      *account.Account
			destination *account.Account
			transfer    *transfer.Transfer
		)
		BeforeEach(func() {
			var err error
			origin, err = account.OpenAccount("origin")
			Expect(err).ToNot(HaveOccurred())
			Expect(origin.DepositMoney(100)).To(Succeed())

			destination, err = account.OpenAccount("destination")
			Expect(err).ToNot(HaveOccurred())

			transfer, err = origin.TransferMoney(50, destination)
			Expect(err).ToNot(HaveOccurred())
		})

		When("it's the origin account", func() {
			It("sends it and decreases the balance", func() {
				err := origin.SendTransfer(transfer)
				Expect(err).ToNot(HaveOccurred())

				Expect(origin.Balance()).To(Equal(50))
			})

			When("the transfer is already sent", func() {
				It("does not send it again", func() {
					err := origin.SendTransfer(transfer)
					Expect(err).ToNot(HaveOccurred())
					Expect(origin.Balance()).To(Equal(50))

					err = origin.SendTransfer(transfer)
					Expect(err).ToNot(HaveOccurred())
					Expect(origin.Balance()).To(Equal(50))
				})
			})
		})

		When("it's the destination account", func() {
			It("assigns to it and increases the balance", func() {
				err := destination.ReceiveTransfer(transfer)
				Expect(err).ToNot(HaveOccurred())

				Expect(destination.Balance()).To(Equal(50))
			})

			When("the transfer is already assigned", func() {
				It("does not assign it again", func() {
					err := destination.ReceiveTransfer(transfer)
					Expect(err).ToNot(HaveOccurred())
					Expect(destination.Balance()).To(Equal(50))

					err = destination.ReceiveTransfer(transfer)
					Expect(err).ToNot(HaveOccurred())
					Expect(destination.Balance()).To(Equal(50))
				})
			})
		})
	})
})
