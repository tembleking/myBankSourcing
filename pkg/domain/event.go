package domain

import "time"

type Event interface {
	AggregateID() string
	EventID() EventID
	EventName() string
	HappenedOn() time.Time
	Version() uint64
}
