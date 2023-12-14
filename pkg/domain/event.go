package domain

import "time"

type Event interface {
	AggregateID() string
	EventID() string
	EventName() string
	HappenedOn() time.Time
	Version() uint64
}
