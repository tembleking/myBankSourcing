package factory

import (
	"context"

	gohttp "net/http"
	"time"

	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tembleking/myBankSourcing/internal/lazy"
	"github.com/tembleking/myBankSourcing/pkg/application/grpc"
	"github.com/tembleking/myBankSourcing/pkg/application/http"
	pb "github.com/tembleking/myBankSourcing/pkg/application/proto"
	"github.com/tembleking/myBankSourcing/pkg/domain"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/domain/views"
	"github.com/tembleking/myBankSourcing/pkg/persistence"
	persistenceaccount "github.com/tembleking/myBankSourcing/pkg/persistence/account"
	"github.com/tembleking/myBankSourcing/pkg/persistence/serializer"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

type Factory struct {
	accountServiceField    lazy.Lazy[*services.AccountService]
	eventStoreField        lazy.Lazy[*persistence.EventStore]
	appendOnlyStoreField   lazy.Lazy[persistence.AppendOnlyStore]
	httpHandlerField       lazy.Lazy[gohttp.Handler]
	grpcServerField        lazy.Lazy[*gogrpc.Server]
	accountViewField       lazy.Lazy[*views.AccountView]
	accountRepositoryField lazy.Lazy[domain.Repository[*account.Account]]
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewAccountService() *services.AccountService {
	return f.accountServiceField.GetOrInit(func() *services.AccountService {
		eventStore := f.eventStore()
		return services.NewAccountService(eventStore, f.accountRepository())
	})
}

func (f *Factory) accountRepository() domain.Repository[*account.Account] {
	return f.accountRepositoryField.GetOrInit(func() domain.Repository[*account.Account] {
		return persistenceaccount.NewRepository(f.eventStore())
	})
}

func (f *Factory) NewAccountView() *views.AccountView {
	return f.accountViewField.GetOrInit(func() *views.AccountView {
		accountView, err := views.NewAccountView(f.eventStore())
		if err != nil {
			panic(err)
		}
		return accountView
	})
}

func (f *Factory) eventStore() *persistence.EventStore {
	return f.eventStoreField.GetOrInit(func() *persistence.EventStore {
		eventSerializer := &serializer.Msgpack{}
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
		return http.NewHTTPServer(ctx, f.NewAccountService(), f.NewAccountView())
	})
}

func (f *Factory) NewGRPCServer() *gogrpc.Server {

	return f.grpcServerField.GetOrInit(func() *gogrpc.Server {
		accountGRPCServer := grpc.NewAccountGRPCServer(f.NewAccountService(), f.NewAccountView())
		grpcServer := gogrpc.NewServer()
		reflection.Register(grpcServer)

		pb.RegisterClerkAPIServiceServer(grpcServer, accountGRPCServer)
		return grpcServer
	})
}
