package account

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type Repository interface {
	SaveAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, accountID ID) (*Account, error)
}
