package factory

import (
	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/inmemory"
	accountpersistence "github.com/tembleking/myBankSourcing/pkg/persistence/inmemory/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
)

func NewAccountService() *services.AccountService {
	eventSerializer := &serializer.Msgpack{}
	eventStore := persistence.NewEventStore(eventSerializer, eventSerializer, inmemory.NewAppendOnlyStore())
	repository := accountpersistence.NewRepository(eventStore)
	return services.NewAccountService(repository)
}
