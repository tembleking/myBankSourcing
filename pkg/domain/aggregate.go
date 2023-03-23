package domain

type BaseAggregate struct {
	events           []Event
	aggregateVersion uint64
	OnEventFunc      func(Event)
}

// Apply applies the event by calling the on-event function and saves them, so they can then be returned by Events
func (a *BaseAggregate) Apply(event Event) {
	a.OnEventFunc(event)
	a.events = append(a.events, event)
	a.aggregateVersion++
}

// Events returns the applied events
func (a *BaseAggregate) Events() []Event {
	events := make([]Event, len(a.events))
	copy(events, a.events)
	return events
}

func (a *BaseAggregate) AggregateVersion() uint64 {
	return a.aggregateVersion
}
