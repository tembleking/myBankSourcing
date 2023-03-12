package domain

import "context"

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type Subscribable interface {
	Subscribe(ctx context.Context) (<-chan Event, error)
}
