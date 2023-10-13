package factory

import (
	"context"
	"fmt"
	gohttp "net/http"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/application/grpc"
	pb "github.com/tembleking/myBankSourcing/pkg/application/proto"
	"github.com/tembleking/myBankSourcing/pkg/broker"
	"github.com/tembleking/myBankSourcing/pkg/outbox"

	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	surreal "github.com/surrealdb/surrealdb.go"

	"github.com/tembleking/myBankSourcing/internal/lazy"
	"github.com/tembleking/myBankSourcing/pkg/application/http"
	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

type Factory struct {
	accountServiceField    lazy.Lazy[*services.AccountService]
	eventStoreField        lazy.Lazy[*persistence.EventStore]
	appendOnlyStoreField   lazy.Lazy[persistence.AppendOnlyStore]
	surrealDBInstanceField lazy.Lazy[*surreal.DB]
	httpHandlerField       lazy.Lazy[gohttp.Handler]
	grpcServerField        lazy.Lazy[*gogrpc.Server]
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewAccountService() *services.AccountService {
	return f.accountServiceField.GetOrInit(func() *services.AccountService {
		eventStore := f.eventStore()
		return services.NewAccountService(eventStore)
	})
}

func (f *Factory) eventStore() *persistence.EventStore {
	return f.eventStoreField.GetOrInit(func() *persistence.EventStore {
		eventSerializer := &serializer.Msgpack{}
		return persistence.NewEventStoreBuilder().
			WithSerializer(eventSerializer).
			WithDeserializer(eventSerializer).
			WithAppendOnlyStore(f.appendOnlyStore()).
			Build()
	})

}

func (f *Factory) appendOnlyStore() persistence.AppendOnlyStore {
	return f.appendOnlyStoreField.GetOrInit(func() persistence.AppendOnlyStore {
		return f.sqliteInstance()
	})
}

func (f *Factory) surrealDBInstance() *surreal.DB {
	return f.surrealDBInstanceField.GetOrInit(func() *surreal.DB {
		db, err := surreal.New("ws://localhost:8000/rpc")
		if err != nil {
			panic(fmt.Errorf("error connecting to surrealdb: %w", err))
		}
		_, err = db.Signin(map[string]string{
			"user": "root",
			"pass": "root",
		})
		if err != nil {
			panic(fmt.Errorf("error signing in to surrealdb: %w", err))
		}
		_, err = db.Use("ns", "db")
		if err != nil {
			panic(fmt.Errorf("error using namespace and database in surrealdb: %w", err))
		}

		return db
	})
}

func (f *Factory) sqliteInstance() *sqlite.AppendOnlyStore {
	appendOnlyStore, err := sqlite.New("file:///tmp/mybankdb.db3")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = appendOnlyStore.MigrateDB(ctx)
	if err != nil {
		panic(err)
	}

	return appendOnlyStore
}

func (f *Factory) NewHTTPHandler(ctx context.Context) gohttp.Handler {
	return f.httpHandlerField.GetOrInit(func() gohttp.Handler {
		return http.NewHTTPServer(ctx, f.NewAccountService())
	})
}

func (f *Factory) NewGRPCServer() *gogrpc.Server {
	return f.grpcServerField.GetOrInit(func() *gogrpc.Server {
		accountGRPCServer := grpc.NewAccountGRPCServer(f.NewAccountService())
		grpcServer := gogrpc.NewServer()
		reflection.Register(grpcServer)

		pb.RegisterClerkAPIServiceServer(grpcServer, accountGRPCServer)
		return grpcServer
	})
}

func (f *Factory) NewTransactionalOutboxPublisher() *outbox.TransactionalOutbox {
	return outbox.NewTransactionalOutboxBuilder(broker.NewInMemoryMessageBroker()).
		WithAppendOnlyStore(f.appendOnlyStore()).
		Build()
}
