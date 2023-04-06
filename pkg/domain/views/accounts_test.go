package views_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/domain/views"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

var _ = Describe("Accounts", func() {
	var (
		eventStore *persistence.EventStore
	)

	BeforeEach(func() {
		eventStore = persistence.NewEventStoreBuilder().Build()

		err := eventStore.AppendToStream(context.Background(), "aggregate-0", 5, []domain.Event{
			&account.AccountOpened{AccountID: "some-account", AccountVersion: 0},
			&account.AmountAdded{Quantity: 50, Balance: 50, AccountVersion: 1},
			&account.AmountAdded{Quantity: 100, Balance: 150, AccountVersion: 2},
			&account.TransferSent{Quantity: 50, Balance: 100, From: "some-account", To: "another-account", AccountVersion: 3},
			&account.AmountWithdrawn{Quantity: 25, Balance: 75, AccountVersion: 4},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	When("there are multiple TransferSent events saved", func() {
		It("returns the transfers found", func() {
			accountsView, err := views.NewAccountView(eventStore)
			Expect(err).ToNot(HaveOccurred())

			transfers := accountsView.Accounts()

			Expect(transfers).To(ConsistOf(
				"some-account",
			))
		})
	})
})
