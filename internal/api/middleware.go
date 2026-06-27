package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/darksuei/kubeRPC/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		slog.Debug("request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery)

		next.ServeHTTP(rec, r)

		metrics.HTTPRequests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rec.status)).Inc()

		level := slog.LevelDebug
		if rec.status >= 400 {
			level = slog.LevelError
		}
		slog.Log(r.Context(), level, "response",
			"status", rec.status,
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})
}
