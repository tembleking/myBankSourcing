package domain

type Event interface {
	AggregateID() string
	EventName() string
	Version() uint64
}
