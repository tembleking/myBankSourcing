package domain

type BaseAggregate struct {
	events      []Event
	OnEventFunc func(Event)
}

// Apply applies the event by calling On and saves them, so they can then be returned by Events
func (a *BaseAggregate) Apply(event Event) {
	if event.ID() == "" {
		panic("ID not assigned to event")
	}
	a.events = append(a.events, event)
	a.On(event)
}

// On executes the modifications of an event to the aggregate. It should only be called by the
// persistence layer.
func (a *BaseAggregate) On(event Event) {
	a.OnEventFunc(event)
}

// Events returns the applied events
func (a *BaseAggregate) Events() []Event {
	return a.events
}
