package serializer

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type GoBinarySerializer struct{}

func (g *GoBinarySerializer) Serialize(events domain.Event) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(&events)
	if err != nil {
		return nil, fmt.Errorf("error serializing events: %w", err)
	}
	return buf.Bytes(), nil
}

func (g *GoBinarySerializer) Deserialize(data []byte) (domain.Event, error) {
	var event domain.Event
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("error deserializing event: %w", err)
	}
	return event, nil
}
