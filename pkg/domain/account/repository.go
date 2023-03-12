package account

import "context"

type Repository interface {
	SaveAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, accountID ID) (*Account, error)
}
