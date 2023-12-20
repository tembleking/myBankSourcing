package transfer

import (
	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Transfer struct {
	domain.BaseAggregate

	TransferID  string
	FromAccount string
	ToAccount   string
	Amount      int
}

func NewTransfer() *Transfer {
	t := &Transfer{}
	t.OnEventFunc = t.onEvent
	return t
}

func RequestTransfer(fromAccount string, toAccount string, amount int) *Transfer {
	transfer := NewTransfer()
	transfer.Apply(&TransferRequested{
		TransferID:  uuid.NewString(),
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
	})
	return transfer
}

func (t *Transfer) SameEntityAs(other domain.Entity) bool {
	if t == nil && other == nil {
		return true
	}
	if otherTransfer, ok := other.(*Transfer); ok {
		return t.TransferID == otherTransfer.TransferID &&
			t.FromAccount == otherTransfer.FromAccount &&
			t.ToAccount == otherTransfer.ToAccount &&
			t.Amount == otherTransfer.Amount
	}
	return false
}

func (t *Transfer) onEvent(event domain.Event) {
	switch e := event.(type) {
	case *TransferRequested:
		t.TransferID = e.TransferID
		t.FromAccount = e.FromAccount
		t.ToAccount = e.ToAccount
		t.Amount = e.Amount
	}
}
