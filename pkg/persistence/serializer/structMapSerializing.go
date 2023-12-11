package serializer

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

type structMapSerializing struct {
	registeredTypes map[string]any
}

var structMapSerializer *structMapSerializing

func init() {
	structMapSerializer = &structMapSerializing{registeredTypes: map[string]any{}}
}

func (s *structMapSerializing) serializeToMap(event domain.Event) (map[string]any, error) {
	result := map[string]any{}
	err := mapstructure.Decode(event, &result)
	if err != nil {
		return nil, fmt.Errorf("error serializing event %s to map: %w", event.EventName(), err)
	}
	return result, nil
}

func (s *structMapSerializing) deserializeFromMap(eventName string, data map[string]any) (domain.Event, error) {
	if _, ok := s.registeredTypes[eventName]; !ok {
		return nil, fmt.Errorf("error deserializing from map, type not registered: %s", eventName)
	}

	event := reflect.New(reflect.TypeOf(s.registeredTypes[eventName])).Interface()
	err := mapstructure.Decode(data, event)
	if err != nil {
		return nil, fmt.Errorf("error deserializing type %s from map: %w", eventName, err)
	}
	return event.(domain.Event), nil
}

func (s *structMapSerializing) register(value domain.Event) {
	reflectValue := reflect.ValueOf(value)
	for {
		if reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			continue
		}
		break
	}
	s.registeredTypes[value.EventName()] = reflectValue.Interface()
}
