package middlewares

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// wrapper type to intercept WriteHeader call so that we can observe
// the statusCode.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "the total number of requests received",
		},
		[]string{"path", "code"},
	)
	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_latency_seconds",
			Help:    "request latency in seconds",
			Buckets: []float64{0.0001, 0.001, 0.01, 0.1},
		},
		[]string{},
	)
)

func Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		requestCounter,
		latencyHistogram,
	}
}

func WithMetrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.EscapedPath()
			rw := newResponseWriter(w)

			start := time.Now()
			next.ServeHTTP(rw, r)
			latency := time.Since(start)
			code := rw.statusCode

			requestCounter.WithLabelValues(path, strconv.Itoa(code)).Inc()
			latencyHistogram.WithLabelValues().Observe(latency.Seconds())
		}

		return http.HandlerFunc(fn)
	}
}

func WithAPIKey(apiKey string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if apiKey != "" {
				slog.Info("Auth middleware: compare Authorization header with target API Key")
				val := r.Header.Get("Authorization")
				if val != apiKey {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Not Authenticated"))
					return
				}
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
