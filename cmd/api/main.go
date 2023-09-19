package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"log/slog"

	"os/signal"
	"syscall"

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
	ENV_SHUTDOWN_TIMEOUT    = "SHUTDOWN_TIMEOUT"
	ENV_MAX_HEADER_BYTES    = "MAX_HEADER_BYTES"
	ENV_API_KEY             = "API_KEY"
	defaultDuration         = time.Second * 3
	defaultShutdownTimeout  = time.Second * 3
	defaultMaxHeaderBytes   = 1024
)

func main() {
	svc := fizzbuzz.NewService()

	// wrap with cancel so that both server an stop eachother
	ctx, cancel := context.WithCancel(
		contextWithSigterm(context.Background()),
	)
	defer cancel()

	var group errgroup.Group

	apiAddr := readAddressFromEnv()
	api := newAPIServer(ctx, apiAddr, svc)
	adminAddr := readAdminAddressFromEnv()
	admin := newAdminServer(ctx, adminAddr)

	group.Go(func() error {
		defer cancel()
		defer slog.Info("admin server stopped")

		slog.Info("start admin server...", "addr", adminAddr)
		return runServer(admin)
	})

	group.Go(func() error {
		defer cancel()
		defer slog.Info("api server stopped")

		slog.Info("start api server...", "addr", apiAddr)
		return runServer(api)
	})

	<-ctx.Done()
	slog.Info("main context canceled")

	timeout := parseDuration(os.Getenv(ENV_SHUTDOWN_TIMEOUT), defaultShutdownTimeout)
	shutDownServer(admin, timeout)
	shutDownServer(api, timeout)

	// we should not need this, but it's good hygiene to wait for the group
	// if we spawn some other goroutines
	err := group.Wait()
	slog.Info("both servers shut down")

	slog.Error("group.Wait()", "err", err)
}

func contextWithSigterm(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		var counter int

		for {
			if counter == 2 {
				slog.Error("force quit!")
				os.Exit(1)
			}

			select {
			case sig := <-signalCh:
				slog.Info("caught signal", "sig", sig)
			case <-ctx.Done():
				slog.Info("context done before signal")
			}

			cancel()
			counter++
		}
	}()

	return ctxWithCancel
}

func runServer(s *http.Server) error {
	return s.ListenAndServe()
}

// a simple admin server for health checks and metrics
func newAdminServer(ctx context.Context, addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &http.Server{}
}

// the fizzbuzz API
func newAPIServer(ctx context.Context, addr string, svc fizzbuzz_v1.FizzBuzzServiceServer) *http.Server {
	// use the grpc gateway runtime as router
	r := runtime.NewServeMux()

	// some toy middlewares
	handler := middlewares.WrapHandler(
		r,
		authn.WithAPIKey(os.Getenv(ENV_API_KEY)),
	)

	// register routes defined in proto annotations
	err := fizzbuzz_v1.RegisterFizzBuzzServiceHandlerServer(ctx, r, svc)
	expect("register fizzbuzz handler", err)

	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		ReadHeaderTimeout: parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		WriteTimeout:      parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		IdleTimeout:       parseDuration(ENV_READ_HEADER_TIMEOUT, defaultDuration),
		MaxHeaderBytes:    parseInt(ENV_MAX_HEADER_BYTES, defaultMaxHeaderBytes),
	}
}

func shutDownServer(s *http.Server, timeout time.Duration) {
	shutCtx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	defer cancel()

	s.Shutdown(shutCtx)
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
