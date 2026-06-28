package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
)

// serviceMetaFields are stored in the same hash as methods but are not methods.
var serviceMetaFields = map[string]bool{
	"serviceName": true,
	"host":        true,
	"port":        true,
}

func GetMethods(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	all, err := cache.Store.HGetAll("service:" + serviceName)
	if err != nil {
		slog.Error("get-methods: cache error", "service", serviceName, "error", err)
		http.Error(w, "Error retrieving service", http.StatusInternalServerError)
		return
	}

	if len(all) == 0 {
		slog.Debug("get-methods: service not found", "service", serviceName)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	methods := make(map[string]json.RawMessage, len(all))
	for k, v := range all {
		if !serviceMetaFields[k] {
			methods[k] = json.RawMessage(v)
		}
	}

	slog.Debug("get-methods: returning", "service", serviceName, "count", len(methods))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}
