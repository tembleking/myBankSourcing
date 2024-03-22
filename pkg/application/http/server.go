package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/tembleking/myBankSourcing/pkg/account"
	"github.com/tembleking/myBankSourcing/pkg/application/grpc"
	"github.com/tembleking/myBankSourcing/pkg/application/proto"
)

func NewHTTPServer(ctx context.Context, accountService *account.Service, accountProjection *account.Projection) http.Handler {
	mux := runtime.NewServeMux()
	err := proto.RegisterClerkAPIServiceHandlerServer(ctx, mux, grpc.NewAccountGRPCServer(accountService, accountProjection))
	if err != nil {
		panic(err)
	}
	return mux
}
