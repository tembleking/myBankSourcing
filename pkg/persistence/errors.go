package persistence

import (
	"fmt"
)

type ErrUnexpectedVersion struct {
	Found    StreamVersion
	Expected StreamVersion
}

func (e *ErrUnexpectedVersion) Error() string {
	return fmt.Sprintf("unexpected version: found %d, expected %d", e.Found, e.Expected)
}
