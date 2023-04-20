package outbox

import (
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

type TransactionalOutboxBuilder struct {
	appendOnlyStore         persistence.AppendOnlyStore
	messageBrokerSerializer MessageBrokerSerializer
	messageBroker           PublishableMessageBroker
}

func NewTransactionalOutboxBuilder(messageBroker PublishableMessageBroker) *TransactionalOutboxBuilder {
	defaultAppendOnlyStore := persistence.NewInMemoryAppendOnlyStore()
	messageBrokerSerializer := &serializer.Msgpack{}

	return &TransactionalOutboxBuilder{
		appendOnlyStore:         defaultAppendOnlyStore,
		messageBroker:           messageBroker,
		messageBrokerSerializer: messageBrokerSerializer,
	}
}

func (b *TransactionalOutboxBuilder) WithAppendOnlyStore(appendOnlyStore persistence.AppendOnlyStore) *TransactionalOutboxBuilder {
	b.appendOnlyStore = appendOnlyStore
	return b
}

func (b *TransactionalOutboxBuilder) WithMessageBrokerPublisher(messageBroker PublishableMessageBroker) *TransactionalOutboxBuilder {
	b.messageBroker = messageBroker
	return b
}

func (b *TransactionalOutboxBuilder) Build() *TransactionalOutbox {
	return &TransactionalOutbox{
		appendOnlyStore:         b.appendOnlyStore,
		messageBroker:           b.messageBroker,
		messageBrokerSerializer: b.messageBrokerSerializer,
	}
}
