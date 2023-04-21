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

		err := eventStore.AppendToStream(context.Background(), "some-account", 5, []domain.Event{
			&account.AccountOpened{AccountID: "some-account", AccountVersion: 0},
			&account.AmountAdded{AccountID: "some-account", Quantity: 50, Balance: 50, AccountVersion: 1},
			&account.AmountAdded{AccountID: "some-account", Quantity: 100, Balance: 150, AccountVersion: 2},
			&account.TransferSent{Quantity: 50, Balance: 100, From: "some-account", To: "another-account", AccountVersion: 3},
			&account.AmountWithdrawn{AccountID: "some-account", Quantity: 25, Balance: 75, AccountVersion: 4},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	When("there are multiple AccountOpened events saved", func() {
		It("returns the transfers found", func() {
			accountsView, err := views.NewAccountView(eventStore)
			Expect(err).ToNot(HaveOccurred())

			accounts := accountsView.Accounts()

			Expect(accounts).To(HaveLen(1))
			Expect(accounts[0].ID()).To(Equal(account.ID("some-account")))
			Expect(accounts[0].Balance()).To(Equal(75))
			Expect(accounts[0].AggregateVersion()).To(Equal(uint64(5)))
		})
	})

	When("a new event is published", func() {
		It("updates the view", func() {
			accountsView, err := views.NewAccountView(eventStore)
			Expect(err).ToNot(HaveOccurred())

			accountsView.Dispatch([]persistence.StreamEvent{
				{ID: persistence.StreamID{StreamName: "another-account", StreamVersion: 0}, Event: &account.AccountOpened{AccountID: "another-account", AccountVersion: 0}},
				{ID: persistence.StreamID{StreamName: "another-account", StreamVersion: 1}, Event: &account.TransferReceived{Quantity: 50, Balance: 50, From: "some-account", To: "another-account", AccountVersion: 1}},
			}...)

			accounts := accountsView.Accounts()

			Expect(accounts).To(HaveLen(2))
			Expect(accounts[0].ID()).To(Equal(account.ID("another-account")))
			Expect(accounts[0].Balance()).To(Equal(50))
			Expect(accounts[1].ID()).To(Equal(account.ID("some-account")))
			Expect(accounts[1].Balance()).To(Equal(75))
			Expect(accounts[1].AggregateVersion()).To(Equal(uint64(5)))
		})
	})
})
