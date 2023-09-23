package surrealdb

import (
	"time"

	surreal "github.com/surrealdb/surrealdb.go"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type surrealStoredStreamEventID struct {
	StreamName    string `json:"stream_name"`
	StreamVersion uint64 `json:"stream_version"`
}

type surrealStoredStreamEventIn struct {
	surrealStoredStreamEventOut
	ID surrealStoredStreamEventID `json:"id"`
}

type surrealStoredStreamEventOut struct {
	surreal.Basemodel `table:"event"`

	EventName     string    `json:"event_name"`
	EventData     []byte    `json:"event_data"`
	HappenedOn    time.Time `json:"happened_on"`
	StreamName    string    `json:"stream_name"`
	StreamVersion uint64    `json:"stream_version"`
}

func (r surrealStoredStreamEventOut) ToStoredStreamEvent() persistence.StoredStreamEvent {
	return persistence.StoredStreamEvent{
		ID: persistence.StreamID{
			StreamName:    persistence.StreamName(r.StreamName),
			StreamVersion: persistence.StreamVersion(r.StreamVersion),
		},
		EventName:  r.EventName,
		EventData:  r.EventData,
		HappenedOn: r.HappenedOn,
	}
}

func storedStreamEventToSurreal(storedStreamEvent persistence.StoredStreamEvent) surrealStoredStreamEventIn {
	return surrealStoredStreamEventIn{
		ID: surrealStoredStreamEventID{
			StreamName:    string(storedStreamEvent.ID.StreamName),
			StreamVersion: uint64(storedStreamEvent.ID.StreamVersion),
		},
		surrealStoredStreamEventOut: surrealStoredStreamEventOut{

			StreamName:    string(storedStreamEvent.ID.StreamName),
			StreamVersion: uint64(storedStreamEvent.ID.StreamVersion),
			EventName:     storedStreamEvent.EventName,
			EventData:     storedStreamEvent.EventData,
			HappenedOn:    storedStreamEvent.HappenedOn,
		},
	}
}
