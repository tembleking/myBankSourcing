package domain

type Event interface {
	EventName() string
	Version() uint64
}
