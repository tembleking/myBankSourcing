package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/tembleking/myBankSourcing/pkg/application/grpc"
	"github.com/tembleking/myBankSourcing/pkg/application/proto"
	"github.com/tembleking/myBankSourcing/pkg/domain/services"
)

func NewHTTPServer(ctx context.Context, accountService *services.AccountService) http.Handler {
	mux := runtime.NewServeMux()
	err := proto.RegisterClerkAPIServiceHandlerServer(ctx, mux, grpc.NewAccountGRPCServer(accountService))
	if err != nil {
		panic(err)
	}
	return mux
}
