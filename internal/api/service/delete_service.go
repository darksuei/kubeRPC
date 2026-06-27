package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
	"github.com/darksuei/kubeRPC/internal/metrics"
)

func DeleteService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	deleted, err := cache.Store.Del("service:" + serviceName)
	if err != nil {
		slog.Error("delete-service: cache error", "service", serviceName, "error", err)
		http.Error(w, "Error deleting service", http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		slog.Debug("delete-service: not found", "service", serviceName)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	metrics.RegisteredServices.Dec()
	slog.Info("service deleted", "service", serviceName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service %s deleted successfully", serviceName)
}
