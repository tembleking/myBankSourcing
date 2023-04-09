package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/tembleking/myBankSourcing/internal/factory"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	wg := &sync.WaitGroup{}
	factory := factory.NewFactory()

	wg.Add(1)
	go serveHTTP(ctx, wg, factory)

	wg.Wait()
	return nil
}

func serveHTTP(ctx context.Context, wg *sync.WaitGroup, factory *factory.Factory) {
	defer wg.Done()

	server := &http.Server{
		Handler: factory.NewHTTPHandler(ctx),
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(fmt.Errorf("error listening HTTP on port 8080: %w", err))
	}
	defer listener.Close()

	fmt.Println("listening on port 8080")
	go func() {
		<-ctx.Done()
		fmt.Println("shutting down HTTP server")
		err := server.Shutdown(ctx)
		if err != nil {
			panic(fmt.Errorf("error shutting down HTTP server: %w", err))
		}
	}()

	err = server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Errorf("error serving HTTP: %w", err))
	}
}
