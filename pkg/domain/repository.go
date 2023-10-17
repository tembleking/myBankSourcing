package domain

import "context"

type Repository[T Aggregate] interface {
	GetByID(ctx context.Context, id string) (T, error)
	Save(ctx context.Context, aggregate T) error
}
