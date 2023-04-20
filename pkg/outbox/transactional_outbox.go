package outbox

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

type AcknowledgableMessage interface {
	Data() []byte
	Ack() error
}

type MessageHandler func(message AcknowledgableMessage)

type SubscribableMessageBroker interface {
	Subscribe(ctx context.Context, f MessageHandler) error
}

//go:generate mockgen -source=$GOFILE -destination=mocks/$GOFILE -package=mocks
type PublishableMessageBroker interface {
	Publish(ctx context.Context, message []byte) error
}

type MessageBrokerSerializer interface {
	Serialize(message map[string]string) ([]byte, error)
}

type TransactionalOutbox struct {
	appendOnlyStore         persistence.AppendOnlyStore
	messageBrokerSerializer MessageBrokerSerializer
	messageBroker           PublishableMessageBroker
}

func (e *TransactionalOutbox) DispatchUndispatchedEvents(ctx context.Context) error {
	records, err := e.appendOnlyStore.ReadUndispatchedRecords(ctx)
	if err != nil {
		return fmt.Errorf("error reading records: %w", err)
	}

	for _, storedStreamEvent := range records {
		messageMap := e.eventToMap(storedStreamEvent)
		data, err := e.messageBrokerSerializer.Serialize(messageMap)
		if err != nil {
			return fmt.Errorf("error serializing message: %w", err)
		}

		err = e.messageBroker.Publish(ctx, data)
		if err != nil {
			return fmt.Errorf("error publishing event: %w", err)
		}

		err = e.appendOnlyStore.MarkRecordsAsDispatched(ctx, storedStreamEvent)
		if err != nil {
			return fmt.Errorf("error marking record as dispatched: %w", err)
		}
	}

	return nil
}

func (e *TransactionalOutbox) eventToMap(event persistence.StoredStreamEvent) map[string]string {
	return map[string]string{
		"stream_id":      event.StreamID,
		"stream_version": fmt.Sprintf("%d", event.StreamVersion),
		"event_name":     event.EventName,
		"event_data":     base64.StdEncoding.EncodeToString(event.EventData),
		"happened_on":    event.HappenedOn.Format(time.RFC3339),
	}
}
