package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/tembleking/myBankSourcing/pkg/application/grpc"
	"github.com/tembleking/myBankSourcing/pkg/application/proto"
	"github.com/tembleking/myBankSourcing/pkg/domain/services"
	"github.com/tembleking/myBankSourcing/pkg/domain/views"
)

func NewHTTPServer(ctx context.Context, accountService *services.AccountService, accountView *views.AccountView) http.Handler {
	mux := runtime.NewServeMux()
	err := proto.RegisterClerkAPIServiceHandlerServer(ctx, mux, grpc.NewAccountGRPCServer(accountService, accountView))
	if err != nil {
		panic(err)
	}
	return mux
}
