package domain

type BaseAggregate struct {
	events      []Event
	OnEventFunc func(Event)
}

// Apply applies the event by calling the on-event function and saves them, so they can then be returned by Events
func (a *BaseAggregate) Apply(event Event) {
	if event.ID() == "" {
		panic("ID not assigned to event")
	}
	a.events = append(a.events, event)
	a.OnEventFunc(event)
}

// Events returns the applied events
func (a *BaseAggregate) Events() []Event {
	events := make([]Event, len(a.events))
	copy(events, a.events)
	return events
}
