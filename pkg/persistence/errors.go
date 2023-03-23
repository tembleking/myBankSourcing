package persistence

import (
	"fmt"
)

type ErrRecordsNotFoundForStream struct {
	StreamID string
}

func (e *ErrRecordsNotFoundForStream) Error() string {
	return fmt.Sprintf("records not found for stream: %s", e.StreamID)
}

type ErrRecordsNotFoundForEvent struct {
	EventName string
}

func (e *ErrRecordsNotFoundForEvent) Error() string {
	return fmt.Sprintf("records not found for event: %s", e.EventName)
}

type ErrUnexpectedVersion struct {
	Found    uint64
	Expected uint64
}

func (e *ErrUnexpectedVersion) Error() string {
	return fmt.Sprintf("unexpected version: found %d, expected %d", e.Found, e.Expected)
}
