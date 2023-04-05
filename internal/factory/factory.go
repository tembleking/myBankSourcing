package factory

import (
	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	accountpersistence "github.com/tembleking/myBankSourcing/pkg/persistence/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func NewAccountService() *services.AccountService {
	eventSerializer := &serializer.Msgpack{}
	eventStore := persistence.NewEventStoreBuilder().
		WithSerializer(eventSerializer).
		WithDeserializer(eventSerializer).
		Build()
	repository := accountpersistence.NewRepository(eventStore)
	return services.NewAccountService(repository)
}
