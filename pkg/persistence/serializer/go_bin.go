package serializer

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

func init() {
	gob.Register(&account.AmountAdded{})
	gob.Register(&account.AmountWithdrawn{})
	gob.Register(&account.TransferenceReceived{})
	gob.Register(&account.TransferenceSent{})
}

type GoBinarySerializer struct{}

func (J *GoBinarySerializer) Serialize(events []domain.Event) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(events)
	if err != nil {
		return nil, fmt.Errorf("error serializing events: %w", err)
	}
	return buf.Bytes(), nil
}

func (J *GoBinarySerializer) Deserialize(data []byte) ([]domain.Event, error) {
	events := make([]domain.Event, 0)
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&events)
	if err != nil {
		return nil, fmt.Errorf("error deserializing events: %w", err)
	}
	return events, nil
}
