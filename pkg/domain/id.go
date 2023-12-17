package domain

import "github.com/google/uuid"

type EventID string

func (e EventID) SameValueObjectAs(other ValueObject) bool {
	if otherEventID, ok := other.(EventID); ok {
		return e == otherEventID
	}
	return false
}

func NewEventID() EventID {
	return EventID(uuid.Must(uuid.NewV7()).String())
}
