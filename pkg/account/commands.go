package account

import "github.com/tembleking/myBankSourcing/pkg/domain"

type OpenNewAccount struct {
	ID string
}

// SameCommandAs implements domain.Command.
func (o *OpenNewAccount) SameCommandAs(other domain.Command) bool {
	otherCommand, ok := other.(*OpenNewAccount)
	return ok && o.ID == otherCommand.ID
}
