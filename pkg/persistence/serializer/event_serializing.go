package serializer

import (
	"encoding/gob"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

func RegisterSerializableEvent(event domain.Event) {
	gob.RegisterName(event.EventName(), event)
	structMapSerializer.register(event)
}
