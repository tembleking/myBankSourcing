package serializer

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type Msgpack struct{}

func (m *Msgpack) Serialize(events []domain.Event) ([]byte, error) {
	eventDataAsMap := make([]map[string]map[string]any, 0)
	for _, event := range events {
		eventData, err := structMapSerializer.SerializeToMap(event)
		if err != nil {
			return nil, fmt.Errorf("error serializing event to map: %w", err)
		}
		eventDataAsMap = append(eventDataAsMap, eventData)
	}

	data, err := msgpack.Marshal(eventDataAsMap)
	if err != nil {
		return nil, fmt.Errorf("error serializing event data map: %w", err)
	}

	return data, nil
}

func (m *Msgpack) Deserialize(data []byte) ([]domain.Event, error) {
	eventDataAsMap := make([]map[string]map[string]any, 0)
	err := msgpack.Unmarshal(data, &eventDataAsMap)
	if err != nil {
		return nil, fmt.Errorf("error deserializing event data map: %w", err)
	}

	events := make([]domain.Event, 0)
	for _, eventData := range eventDataAsMap {
		event, err := structMapSerializer.DeserializeFromMap(eventData)
		if err != nil {
			return nil, fmt.Errorf("error deserializing event from map: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
