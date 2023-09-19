package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lahssenk/fizzbuzz-api/middlewares"
	"github.com/lahssenk/fizzbuzz-api/middlewares/authn"
	fizzbuzz_v1 "github.com/lahssenk/fizzbuzz-api/protogen/fizzbuzz/v1"
	"github.com/lahssenk/fizzbuzz-api/service/fizzbuzz"
	"golang.org/x/sync/errgroup"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ENV_SERVER_HOST         = "SERVER_HOST"
	ENV_SERVER_PORT         = "SERVER_PORT"
	ENV_ADMIN_PORT          = "ADMIN_PORT"
	ENV_READ_TIMEOUT        = "READ_TIMEOUT"
	ENV_READ_HEADER_TIMEOUT = "READ_HEADER_TIMEOUT"
	ENV_WRITE_TIMEOUT       = "WRITE_TIMEOUT"
	ENV_IDLE_TIMEOUT        = "IDLE_TIMEOUT"
	ENV_MAX_HEADER_BYTES    = "MAX_HEADER_BYTES"
	ENV_API_KEY             = "API_KEY"
	defaultDuration         = time.Second * 3
	defaultMaxHeaderBytes   = 1024
)

func main() {
	svc := fizzbuzz.NewService()
	ctx := context.Background()

	var group errgroup.Group

	group.Go(func() error {
		addr := readAddressFromEnv()
		return runServer(ctx, addr, svc)
	})

	group.Go(func() error {
		addr := readAdminAddressFromEnv()
		return runAdminServer(ctx, addr)
	})

	err := group.Wait()
	expect("run servers", err)
}

func runAdminServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return http.ListenAndServe(addr, mux)
}

func runServer(ctx context.Context, addr string, svc fizzbuzz_v1.FizzBuzzServiceServer) error {
	// use the grpc gateway runtime as router
	r := runtime.NewServeMux()

	handler := middlewares.WrapHandler(
		r,
		authn.WithAPIKey(os.Getenv(ENV_API_KEY)),
	)

	// register routes defined in proto annotations
	err := fizzbuzz_v1.RegisterFizzBuzzServiceHandlerServer(ctx, r, svc)
	expect("register fizzbuzz handler", err)

	server := http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		ReadHeaderTimeout: parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		WriteTimeout:      parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		IdleTimeout:       parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		MaxHeaderBytes:    parseInt(ENV_MAX_HEADER_BYTES, defaultMaxHeaderBytes),
	}

	slog.Info("start fizzbuzz server...", "addr", addr)

	err = server.ListenAndServe()
	expect("listen and server", err)

	return nil
}

func expect(op string, err error) {
	if err != nil {
		slog.Error(op, "err", err)
		os.Exit(1)
	}
}

func readAddressFromEnv() string {
	host := os.Getenv(ENV_SERVER_HOST)
	port := os.Getenv(ENV_SERVER_PORT)

	addr := net.JoinHostPort(host, port)

	if addr == ":" {
		slog.Error("invalid addr ':'")
		os.Exit(1)
	}

	return addr
}

func readAdminAddressFromEnv() string {
	host := os.Getenv(ENV_SERVER_HOST)
	port := os.Getenv(ENV_ADMIN_PORT)

	addr := net.JoinHostPort(host, port)

	if addr == ":" {
		slog.Error("invalid admin addr ':'")
		os.Exit(1)
	}

	return addr
}

func parseDuration(s string, fallback time.Duration) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		dur = fallback
	}

	return dur
}

func parseInt(s string, fallback int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		i = fallback
	}

	return i
}
