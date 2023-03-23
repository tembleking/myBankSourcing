package persistence

import (
	"fmt"
)

type ErrRecordsNotFound struct {
	StreamID string
}

func (e *ErrRecordsNotFound) Error() string {
	return fmt.Sprintf("records not found for stream: %s", e.StreamID)
}

type ErrUnexpectedVersion struct {
	Found    uint64
	Expected uint64
}

func (e *ErrUnexpectedVersion) Error() string {
	return fmt.Sprintf("unexpected version: found %d, expected %d", e.Found, e.Expected)
}
