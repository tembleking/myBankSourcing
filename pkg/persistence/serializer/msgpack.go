package serializer

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Msgpack struct{}

func (m *Msgpack) SerializeDomainEvent(event domain.Event) ([]byte, error) {
	eventData, err := structMapSerializer.SerializeToMap(event)
	if err != nil {
		return nil, fmt.Errorf("error serializing event to map: %w", err)
	}

	data, err := msgpack.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("error serializing event data map: %w", err)
	}

	return data, nil
}

func (m *Msgpack) DeserializeDomainEvent(data []byte) (domain.Event, error) {
	eventDataAsMap := make(map[string]map[string]any)
	err := msgpack.Unmarshal(data, &eventDataAsMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing event data map: %w", err)
	}

	event, err := structMapSerializer.DeserializeFromMap(eventDataAsMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing event from map: %w", err)
	}

	return event, nil
}

func (m *Msgpack) Serialize(event map[string]string) ([]byte, error) {
	data, err := msgpack.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error serializing map: %w", err)
	}

	return data, nil
}
