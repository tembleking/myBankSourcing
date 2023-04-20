package broker

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tembleking/myBankSourcing/pkg/outbox"
)

type inMemoryStoredMessage struct {
	data         []byte
	isAcked      bool
	isProcessing bool
}

type inMemoryAcknowledgableMessage struct {
	*inMemoryStoredMessage
	ack func() error
}

func (a *inMemoryAcknowledgableMessage) Data() []byte {
	return a.data
}

func (a *inMemoryAcknowledgableMessage) Ack() error {
	return a.ack()
}

type InMemoryMessageBroker struct {
	messages      map[string][]*inMemoryStoredMessage
	rwMutex       sync.RWMutex
	subscribers   map[string]outbox.MessageHandler
	servingTicker *time.Ticker
	cleanupTicker *time.Ticker
}

func (i *InMemoryMessageBroker) Publish(_ context.Context, data []byte) error {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()

	for subscriber := range i.subscribers {
		i.messages[subscriber] = append(i.messages[subscriber], &inMemoryStoredMessage{data: data})
	}

	return nil
}

func (i *InMemoryMessageBroker) StartServing(ctx context.Context) {
	defer i.servingTicker.Stop()
	defer i.cleanupTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-i.servingTicker.C:
			i.serveMessages()
		case <-i.cleanupTicker.C:
			i.cleanupAckedMessages()
		}
	}
}

func (i *InMemoryMessageBroker) serveMessages() {
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()

	for subscriber, messages := range i.messages {
		messagesNotAckd := messagesNotACKd(messages)
		messagesToProcess := messagesNotBeingProcessed(messagesNotAckd)
		go i.serveMessagesForSubscriber(subscriber, messagesToProcess)
	}
}

func messagesNotBeingProcessed(messagesNotAckd []*inMemoryStoredMessage) []*inMemoryStoredMessage {
	messagesToProcess := make([]*inMemoryStoredMessage, 0, len(messagesNotAckd))
	for _, message := range messagesNotAckd {
		if message.isProcessing {
			continue
		}
		messagesToProcess = append(messagesToProcess, message)
	}
	return messagesToProcess
}

func messagesNotACKd(messages []*inMemoryStoredMessage) []*inMemoryStoredMessage {
	messagesNotAckd := make([]*inMemoryStoredMessage, 0, len(messages))
	for _, message := range messages {
		if message.isAcked {
			continue
		}
		messagesNotAckd = append(messagesNotAckd, message)
	}
	return messagesNotAckd
}

func (i *InMemoryMessageBroker) Subscribe(ctx context.Context, f outbox.MessageHandler) error {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()

	i.subscribers[uuid.NewString()] = f
	return nil
}

func (i *InMemoryMessageBroker) serveMessagesForSubscriber(subscriber string, messagesToProcess []*inMemoryStoredMessage) {
	for _, message := range messagesToProcess {
		i.subscribers[subscriber](&inMemoryAcknowledgableMessage{
			inMemoryStoredMessage: message,
			ack: func() error {
				message.isAcked = true
				return nil
			},
		})
		message.isProcessing = false
	}
}

func (i *InMemoryMessageBroker) cleanupAckedMessages() {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()

	for subscriber, messages := range i.messages {
		nonAckedMessages := messagesNotACKd(messages)
		i.messages[subscriber] = nonAckedMessages
	}
}

func NewInMemoryMessageBroker() *InMemoryMessageBroker {
	return &InMemoryMessageBroker{
		messages:      map[string][]*inMemoryStoredMessage{},
		subscribers:   map[string]outbox.MessageHandler{},
		servingTicker: time.NewTicker(100 * time.Millisecond),
		cleanupTicker: time.NewTicker(1 * time.Second),
	}
}
