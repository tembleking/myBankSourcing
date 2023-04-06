package factory

import (
	surreal "github.com/surrealdb/surrealdb.go"

	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	accountpersistence "github.com/tembleking/myBankSourcing/pkg/persistence/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	"github.com/tembleking/myBankSourcing/pkg/persistence/surrealdb"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewAccountService() *services.AccountService {
	repository := accountpersistence.NewRepository(f.eventStore())
	return services.NewAccountService(repository)
}

func (f *Factory) eventStore() *persistence.EventStore {
	eventSerializer := &serializer.Msgpack{}
	return persistence.NewEventStoreBuilder().
		WithSerializer(eventSerializer).
		WithDeserializer(eventSerializer).
		WithAppendOnlyStore(f.appendOnlyStore()).
		Build()
}

func (f *Factory) appendOnlyStore() persistence.AppendOnlyStore {
	return surrealdb.NewAppendOnlyStore(f.surrealDBInstance())
}

func (f *Factory) surrealDBInstance() *surreal.DB {
	db, err := surreal.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}
	_, err = db.Signin(map[string]string{
		"user": "root",
		"pass": "root",
	})
	if err != nil {
		panic(err)
	}
	_, err = db.Use("ns", "db")
	if err != nil {
		panic(err)
	}
	return db
}
