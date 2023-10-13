package persistence

import (
	"fmt"
)

type ErrUnexpectedVersion struct {
	StreamName StreamName
	Expected   StreamVersion
}

func (e *ErrUnexpectedVersion) Error() string {
	return fmt.Sprintf("unexpected version for stream %s with version %d", e.StreamName, e.Expected)
}
