package serializer

import (
	"encoding/gob"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type EventSerializer interface {
	Serialize(events domain.Event) ([]byte, error)
}

type EventDeserializer interface {
	Deserialize(data []byte) (domain.Event, error)
}

func RegisterSerializableType(t any) {
	gob.Register(t)
	registerTypeInStructSerializer(t)
}
