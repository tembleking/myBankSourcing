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

func registerTypeInStructSerializer(value any) {
	structMapSerializer.Register(value)
}

func (s *structMapSerializing) SerializeToMap(event domain.Event) (map[string]map[string]any, error) {
	name := nameOfValue(event)
	result := map[string]any{}
	err := mapstructure.Decode(event, &result)
	if err != nil {
		return nil, fmt.Errorf("error serializing event %s to map: %w", name, err)
	}
	return map[string]map[string]any{
		name: result,
	}, nil
}

func (s *structMapSerializing) DeserializeFromMap(data map[string]map[string]any) (domain.Event, error) {
	for name, data := range data {
		if _, ok := s.registeredTypes[name]; !ok {
			return nil, fmt.Errorf("error deserializing from map, type not registered: %s", name)
		}

		event := reflect.New(reflect.TypeOf(s.registeredTypes[name])).Interface()
		err := mapstructure.Decode(data, event)
		if err != nil {
			return nil, fmt.Errorf("error deserializing type %s from map: %w", name, err)
		}
		return event.(domain.Event), nil
	}
	return nil, fmt.Errorf("error deserializing event from map: no event found")
}

func (s *structMapSerializing) Register(value any) {
	name := nameOfValue(value)

	s.registerName(name, value)
}

func nameOfValue(value any) string {
	rt := reflect.TypeOf(value)
	name := rt.String()

	if !isStruct(reflect.ValueOf(value)) {
		panic(fmt.Sprintf("error registering type %s, type is not a struct", name))
	}

	star := ""
	if rt.Name() == "" {
		if pt := rt; pt.Kind() == reflect.Pointer {
			star = "*"
			rt = pt
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			name = star + rt.Name()
		} else {
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}
	return name
}

func isStruct(rt reflect.Value) bool {
	if rt.Kind() == reflect.Ptr {
		return isStruct(rt.Elem())
	}
	return rt.Kind() == reflect.Struct
}

func (s *structMapSerializing) registerName(name string, value any) {
	reflectValue := reflect.ValueOf(value)
	for {
		if reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			continue
		}
		break
	}
	s.registeredTypes[name] = reflectValue.Interface()
}
