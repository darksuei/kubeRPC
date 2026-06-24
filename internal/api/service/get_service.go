package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
)

func GetService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	serviceDetails, err := cache.Store.HGetAll("service:" + serviceName)
	if err != nil {
		slog.Error("get-service: cache error", "service", serviceName, "error", err)
		http.Error(w, "Error retrieving service", http.StatusInternalServerError)
		return
	}

	if len(serviceDetails) == 0 {
		slog.Debug("get-service: not found", "service", serviceName)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	slog.Debug("get-service: found", "service", serviceName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serviceDetails)
}
