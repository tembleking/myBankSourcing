package persistence

import (
	"fmt"
)

type ErrUnexpectedVersion struct {
	Found    uint64
	Expected uint64
}

func (e *ErrUnexpectedVersion) Error() string {
	return fmt.Sprintf("unexpected version: found %d, expected %d", e.Found, e.Expected)
}
