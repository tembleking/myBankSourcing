package serializer

import (
	"encoding/gob"
)

func RegisterSerializableType(t any) {
	gob.Register(t)
	registerTypeInStructSerializer(t)
}
