package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/lahssenk/fizzbuzz-api/pkg/api"
	"github.com/lahssenk/fizzbuzz-api/pkg/middlewares"
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
	apiHandlers := api.NewAPI()
	// wrap with cancel so that both servers can stop eachother
	ctx, cancel := context.WithCancel(
		contextWithSigterm(context.Background()),
	)
	defer cancel()

	var group errgroup.Group

	// admin and API server are separate because we probably don't want
	// to expose metrics etc to our consumers
	apiAddr := readAddressFromEnv(ENV_SERVER_HOST, ENV_SERVER_PORT)
	apiServer := newAPIServer(ctx, apiAddr, apiHandlers)
	adminAddr := readAddressFromEnv(ENV_SERVER_HOST, ENV_ADMIN_PORT)
	adminServer := newAdminServer(ctx, adminAddr)

	// spawn a goroutine for each server
	group.Go(func() error {
		defer cancel()
		defer slog.Info("admin server stopped")

		slog.Info("start admin server...", "addr", adminAddr)
		return runServer(adminServer)
	})

	group.Go(func() error {
		defer cancel()
		defer slog.Info("api server stopped")

		slog.Info("start api server...", "addr", apiAddr)
		return runServer(apiServer)
	})

	// wait until ctrl+c or one of the server failed
	<-ctx.Done()
	slog.Info("main context canceled")

	// shutdown remaining server(s)
	timeout := parseDuration(os.Getenv(ENV_SHUTDOWN_TIMEOUT), defaultShutdownTimeout)
	shutDownServer(adminServer, timeout)
	shutDownServer(apiServer, timeout)

	// we should not need this, but it's good hygiene to wait for the group
	// if we spawn some other goroutines
	err := group.Wait()

	slog.Info("both servers shut down")

	expect("group.Wait()", err)
}

func contextWithSigterm(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		// register a chan to receive notification of os signals
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		var counter int

		// we want to propagate the signal for graceful shutdown, but also
		// have a force quit trigger if signal received twice
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

// A simple admin server for health checks and metrics.
// This is listening on another port because admin endpoint
// should not be access by the consumer of the API
func newAdminServer(ctx context.Context, addr string) *http.Server {
	// use a custom registry to remove the default go runtime metrics
	// from promhttp
	registry := prometheus.NewRegistry()
	registry.MustRegister(middlewares.Metrics()...)
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})

	r := chi.NewRouter()
	r.Use(
		middlewares.WithLogger(),
	)

	r.Get(
		"/health",
		func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) },
	)
	r.Get("/metrics", promHandler.ServeHTTP)

	return prepareServer(addr, r)
}

func prepareServer(addr string, handler http.Handler) *http.Server {
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

// The fizzbuzz API server.
func newAPIServer(
	ctx context.Context,
	addr string,
	handlers api.API,
) *http.Server {
	r := chi.NewRouter()
	r.Use(
		// this middleware will info about each request
		middlewares.WithLogger(),
		// this middleware will collect latency and count metrics
		middlewares.WithMetrics(),
		// this middleware will ensure the Authorization header matches the target APIKey if not empty
		middlewares.WithAPIKey(os.Getenv(ENV_API_KEY)),
	)

	r.Get("/fizzbuzz", handlers.FizzBuzzHandler)

	return prepareServer(addr, r)
}

func shutDownServer(s *http.Server, timeout time.Duration) {
	shutCtx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	defer cancel()

	s.Shutdown(shutCtx)
}

// simple util
func expect(op string, err error) {
	if err != nil {
		slog.Error(op, "err", err)
		os.Exit(1)
	}
}

func readAddressFromEnv(hostVarName, portVarName string) string {
	host := os.Getenv(hostVarName)
	port := os.Getenv(portVarName)

	addr := net.JoinHostPort(host, port)

	if addr == ":" {
		slog.Error("invalid addr ':'")
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
