package transferences_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/domain/mocks"
	"github.com/tembleking/myBankSourcing/pkg/domain/views/transferences"
)

var _ = Describe("Transferences", func() {
	var (
		ctrl       *gomock.Controller
		repository *mocks.MockSubscribable
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repository = mocks.NewMockSubscribable(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("there is a transference", func() {
		It("is registered in the view", func() {
			events := make(chan domain.Event)
			repository.EXPECT().Subscribe(gomock.Any()).Return(events, nil)

			view, err := transferences.NewViewSubscribedTo(context.Background(), repository)
			Expect(err).ToNot(HaveOccurred())

			events <- &account.TransferenceSent{
				Quantity: 50,
				From:     "origin",
				To:       "destination",
			}

			Expect(view.Transferences()).To(Equal([]transferences.Transference{{
				Origin:      "origin",
				Destination: "destination",
				Quantity:    50,
			}}))
		})
	})
})
