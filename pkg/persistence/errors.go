package persistence

import (
	"errors"
)

var ErrUnexpectedVersion = errors.New("unexpected version for stream")
