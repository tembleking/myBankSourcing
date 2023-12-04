package domain

import (
	"time"
)

// TODO convert the aggregate into an interface

type Aggregate interface {
	Entity
	UncommittedEvents() []Event
	LoadFromHistory(events ...Event)
	Version() uint64
}

type BaseAggregate struct {
	id          string
	events      []Event
	version     uint64
	OnEventFunc func(Event)
}

// Apply applies the event by calling the on-event function and saves them, so they can then be returned by Events
func (a *BaseAggregate) Apply(event Event) {
	a.apply(event, true)
}

func (a *BaseAggregate) LoadFromHistory(events ...Event) {
	for _, event := range events {
		a.apply(event, false)
	}
}

func (a *BaseAggregate) apply(event Event, isNew bool) {
	a.updateMetadata(event)
	if a.OnEventFunc != nil {
		a.OnEventFunc(event)
	}

	if isNew {
		a.events = append(a.events, event)
	}
}

func (a *BaseAggregate) ID() string {
	return a.id
}

func (a *BaseAggregate) updateMetadata(event Event) {
	a.id = event.AggregateID()
	a.version = event.Version()
}

func (a *BaseAggregate) UncommittedEvents() []Event {
	return a.events
}

func (a *BaseAggregate) Version() uint64 {
	return a.version
}

func (a *BaseAggregate) Now() time.Time {
	return time.Now().UTC()
}

func (a *BaseAggregate) NextVersion() uint64 {
	return a.version + 1
}
