package factory

import (
	"context"
	gohttp "net/http"
	"time"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/transfer"

	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tembleking/myBankSourcing/internal/lazy"
	"github.com/tembleking/myBankSourcing/pkg/application/grpc"
	"github.com/tembleking/myBankSourcing/pkg/application/http"
	pb "github.com/tembleking/myBankSourcing/pkg/application/proto"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

type Factory struct {
	accountServiceField     lazy.Lazy[*account.Service]
	eventStoreField         lazy.Lazy[*persistence.EventStore]
	appendOnlyStoreField    lazy.Lazy[persistence.AppendOnlyStore]
	httpHandlerField        lazy.Lazy[gohttp.Handler]
	grpcServerField         lazy.Lazy[*gogrpc.Server]
	accountProjectionField  lazy.Lazy[*account.Projection]
	accountRepositoryField  lazy.Lazy[domain.Repository[*account.Account]]
	transferRepositoryField lazy.Lazy[domain.Repository[*transfer.Transfer]]
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewAccountService() *account.Service {
	return f.accountServiceField.GetOrInit(func() *account.Service {
		return account.NewAccountService(f.accountRepository(), f.transferRepository())
	})
}

func (f *Factory) accountRepository() domain.Repository[*account.Account] {
	return f.accountRepositoryField.GetOrInit(func() domain.Repository[*account.Account] {
		return account.NewRepository(f.eventStore())
	})
}

func (f *Factory) transferRepository() domain.Repository[*transfer.Transfer] {
	return f.transferRepositoryField.GetOrInit(func() domain.Repository[*transfer.Transfer] {
		return transfer.NewRepository(f.eventStore())
	})
}

func (f *Factory) NewAccountProjection(ctx context.Context) *account.Projection {
	return f.accountProjectionField.GetOrInit(func() *account.Projection {
		accountProjection, err := account.NewAccountProjection(ctx, f.eventStore().ReadOnlyEventStore, time.Second)
		if err != nil {
			panic(err)
		}
		return accountProjection
	})
}

func (f *Factory) eventStore() *persistence.EventStore {
	return f.eventStoreField.GetOrInit(func() *persistence.EventStore {
		eventSerializer := &serializer.JSON{}
		return persistence.NewEventStoreBuilder(f.appendOnlyStore()).
			WithSerializer(eventSerializer).
			WithDeserializer(eventSerializer).
			Build()
	})
}

func (f *Factory) appendOnlyStore() persistence.AppendOnlyStore {
	return f.appendOnlyStoreField.GetOrInit(func() persistence.AppendOnlyStore {
		return f.sqliteInstance()
	})
}

func (f *Factory) sqliteInstance() *sqlite.AppendOnlyStore {
	appendOnlyStore, err := sqlite.New("file:///tmp/mybankdb.sqlite")
	if err != nil {
		panic(err)
	}

	err = appendOnlyStore.MigrateDB()
	if err != nil {
		panic(err)
	}

	return appendOnlyStore
}

func (f *Factory) NewHTTPHandler(ctx context.Context) gohttp.Handler {
	return f.httpHandlerField.GetOrInit(func() gohttp.Handler {
		return http.NewHTTPServer(ctx, f.NewAccountService(), f.NewAccountProjection(ctx))
	})
}

func (f *Factory) NewGRPCServer(ctx context.Context) *gogrpc.Server {
	return f.grpcServerField.GetOrInit(func() *gogrpc.Server {
		accountGRPCServer := grpc.NewAccountGRPCServer(f.NewAccountService(), f.NewAccountProjection(ctx))
		grpcServer := gogrpc.NewServer()
		reflection.Register(grpcServer)

		pb.RegisterClerkAPIServiceServer(grpcServer, accountGRPCServer)
		return grpcServer
	})
}
