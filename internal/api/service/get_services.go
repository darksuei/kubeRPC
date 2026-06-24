package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
)

func GetServices(w http.ResponseWriter, _ *http.Request) {
	keys, err := cache.Store.Keys("service:*")
	if err != nil {
		slog.Error("get-services: cache error", "error", err)
		http.Error(w, "Error retrieving services", http.StatusInternalServerError)
		return
	}

	if len(keys) == 0 {
		slog.Debug("get-services: no services registered")
		http.Error(w, "No services found", http.StatusNotFound)
		return
	}

	services := make(map[string]interface{})
	for _, key := range keys {
		details, err := cache.Store.HGetAll(key)
		if err != nil {
			slog.Error("get-services: error reading key", "key", key, "error", err)
			http.Error(w, "Error retrieving service", http.StatusInternalServerError)
			return
		}
		services[key] = details
	}

	slog.Debug("get-services: returning", "count", len(services))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}
