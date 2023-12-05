package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type JSON struct{}

func (m *JSON) SerializeDomainEvent(event domain.Event) ([]byte, error) {
	eventData, err := structMapSerializer.SerializeToMap(event)
	if err != nil {
		return nil, fmt.Errorf("error serializing event to map: %w", err)
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("error serializing event data map: %w", err)
	}

	return data, nil
}

func (m *JSON) DeserializeDomainEvent(eventName string, data []byte) (domain.Event, error) {
	eventDataAsMap := make(map[string]any)
	err := json.Unmarshal(data, &eventDataAsMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing event data map: %w", err)
	}

	event, err := structMapSerializer.DeserializeFromMap(eventName, eventDataAsMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing event from map: %w", err)
	}

	return event, nil
}
