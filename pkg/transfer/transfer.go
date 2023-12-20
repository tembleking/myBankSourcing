package transfer

import (
	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Transfer struct {
	domain.BaseAggregate

	transferID  string
	fromAccount string
	toAccount   string
	amount      int
}

func (t *Transfer) FromAccount() string {
	return t.fromAccount
}

func (t *Transfer) ToAccount() string {
	return t.toAccount
}

func (t *Transfer) Amount() int {
	return t.amount
}

func (t *Transfer) ID() string {
	return t.transferID
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
		return t.transferID == otherTransfer.transferID &&
			t.fromAccount == otherTransfer.fromAccount &&
			t.toAccount == otherTransfer.toAccount &&
			t.amount == otherTransfer.amount
	}
	return false
}

func (t *Transfer) onEvent(event domain.Event) {
	switch e := event.(type) {
	case *TransferRequested:
		t.transferID = e.TransferID
		t.fromAccount = e.FromAccount
		t.toAccount = e.ToAccount
		t.amount = e.Amount
	}
}
