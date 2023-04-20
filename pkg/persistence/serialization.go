package persistence

import (
	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type DomainEventSerializer interface {
	SerializeDomainEvent(events domain.Event) ([]byte, error)
}

type DomainEventDeserializer interface {
	DeserializeDomainEvent(data []byte) (domain.Event, error)
}
