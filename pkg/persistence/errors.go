package persistence

import (
	"fmt"
)

type ErrAggregateNotFound struct {
	Name string
}

func (e *ErrAggregateNotFound) Error() string {
	return fmt.Sprintf("aggregate not found: %s", e.Name)
}

type ErrUnexpectedVersion struct {
	Found    uint64
	Expected uint64
}

func (e *ErrUnexpectedVersion) Error() string {
	return fmt.Sprintf("unexpected version: found %d, expected %d", e.Found, e.Expected)
}
