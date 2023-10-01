package domain

// TODO convert the aggregate into an interface

type Aggregate interface {
	ID() string
	Events() []Event
	Version() uint64
}

type BaseAggregate struct {
	events      []Event
	version     uint64
	OnEventFunc func(Event)
}

// Apply applies the event by calling the on-event function and saves them, so they can then be returned by Events
func (a *BaseAggregate) Apply(event Event) {
	if event.Version() != a.Version() {
		return
	}

	a.OnEventFunc(event)
	a.events = append(a.events, event)
	a.version++
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
