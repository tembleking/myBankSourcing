package domain

import "fmt"

// TODO convert the aggregate into an interface

type Aggregate interface {
	ID() string
	Events() []Event
	ClearEvents()
	Version() uint64
}

type BaseAggregate struct {
	events      []Event
	version     uint64
	OnEventFunc func(Event)
}

// Apply applies the event by calling the on-event function and saves them, so they can then be returned by Events
func (a *BaseAggregate) Apply(event Event) error {
	if event.Version() != a.Version() {
		return fmt.Errorf("event version '%d' does not match aggregate version '%d'", event.Version(), a.Version())
	}

	a.OnEventFunc(event)
	a.events = append(a.events, event)
	a.version++

	return nil
}

// Events returns the applied events
func (a *BaseAggregate) Events() []Event {
	events := make([]Event, len(a.events))
	copy(events, a.events)
	return events
}

func (a *BaseAggregate) Version() uint64 {
	return a.version
}

func (a *BaseAggregate) ClearEvents() {
	a.events = nil
}
