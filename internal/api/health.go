package api

import (
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if err := cache.Store.Ping(); err != nil {
		slog.Error("health check failed", "error", err)
		http.Error(w, "Cache connection failed", http.StatusInternalServerError)
		return
	}
	slog.Debug("health check passed")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("kubeRPC core is healthy!"))
}
