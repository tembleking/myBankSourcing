package inmemory

import (
	"context"
	"errors"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Repository[T domain.Aggregate] struct {
	aggregates map[string]T
}

func NewRepository[T domain.Aggregate]() *Repository[T] {
	return &Repository[T]{aggregates: make(map[string]T)}
}

func (r *Repository[T]) NextID() string {
	return domain.NewUUID()
}

func (r *Repository[T]) GetByID(_ context.Context, id string) (T, error) {
	if aggregate, ok := r.aggregates[id]; ok {
		return aggregate, nil
	}
	var zero T
	return zero, errors.New("aggregate not found")
}

func (r *Repository[T]) Save(_ context.Context, aggregate T) error {
	r.aggregates[aggregate.ID()] = aggregate
	return nil
}
