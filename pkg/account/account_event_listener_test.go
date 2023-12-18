package account_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/domain"
)

var _ = Describe("AccountEventListener", func() {
	var accountEventListener *account.EventListener
	var accountRepository domain.Repository[*account.Account]

	BeforeEach(func(ctx context.Context) {
		accountRepository = &fakeAccountRepository{map[string]*account.Account{}}
		originAccount, err := account.OpenAccount("origin-account")
		Expect(err).ToNot(HaveOccurred())
		Expect(originAccount.DepositMoney(200)).To(Succeed())
		Expect(accountRepository.Save(ctx, originAccount)).To(Succeed())

		accountEventListener = account.NewEventListener(accountRepository)
	})

	When("a transfer is requested", func() {
		When("and the destination account exists", func() {
			BeforeEach(func() {
				destinationAccount, err := account.OpenAccount("destination-account")
				Expect(err).ToNot(HaveOccurred())
				Expect(accountRepository.Save(context.Background(), destinationAccount)).To(Succeed())
			})

			It("transfers money from one account to another", func(ctx context.Context) {
				err := accountEventListener.OnEvent(ctx, &account.TransferRequested{
					ID:             domain.EventID("some-event-id"),
					TransferID:     "some-transfer-id",
					Quantity:       50,
					Balance:        200,
					From:           "origin-account",
					To:             "destination-account",
					AccountVersion: 0,
					Timestamp:      time.Time{},
				})
				Expect(err).ToNot(HaveOccurred())

				originAccount, err := accountRepository.GetByID(ctx, "origin-account")
				Expect(err).ToNot(HaveOccurred())
				Expect(originAccount.Balance()).To(Equal(200))

				destinationAccount, err := accountRepository.GetByID(ctx, "destination-account")
				Expect(err).ToNot(HaveOccurred())
				Expect(destinationAccount.Balance()).To(Equal(50))
			})
		})

		When("and the destination account does not exist", func() {
			It("rejects the transfer and returns the money to the origin account", func(ctx context.Context) {
				err := accountEventListener.OnEvent(ctx, &account.TransferRequested{
					ID:             domain.EventID("some-event-id"),
					TransferID:     "some-transfer-id",
					Quantity:       50,
					Balance:        200,
					From:           "origin-account",
					To:             "destination-account",
					AccountVersion: 0,
					Timestamp:      time.Time{},
				})
				Expect(err).ToNot(HaveOccurred())

				originAccount, err := accountRepository.GetByID(ctx, "origin-account")
				Expect(err).ToNot(HaveOccurred())
				Expect(originAccount.Balance()).To(Equal(250))
			})
		})
	})
})
