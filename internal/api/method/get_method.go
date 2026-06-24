package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
)

func GetMethod(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	methodName := r.URL.Query().Get("method")
	if methodName == "" {
		http.Error(w, "Method name is required", http.StatusBadRequest)
		return
	}

	slog.Debug("get-method: lookup", "service", serviceName, "method", methodName)

	methodJSON, err := cache.Store.HGet("service:"+serviceName, methodName)
	if err != nil {
		if err == cache.ErrNotFound {
			slog.Error("get-method: not found", "service", serviceName, "method", methodName)
			http.Error(w, "Service or method not found", http.StatusNotFound)
			return
		}
		slog.Error("get-method: cache error", "service", serviceName, "method", methodName, "error", err)
		http.Error(w, "Error retrieving method", http.StatusInternalServerError)
		return
	}

	host, err := cache.Store.HGet("service:"+serviceName, "host")
	if err != nil {
		slog.Error("get-method: host not found", "service", serviceName)
		http.Error(w, "Service host not registered", http.StatusInternalServerError)
		return
	}

	port, err := cache.Store.HGet("service:"+serviceName, "port")
	if err != nil {
		slog.Error("get-method: port not found", "service", serviceName)
		http.Error(w, "Service port not registered", http.StatusInternalServerError)
		return
	}

	slog.Info("method resolved", "service", serviceName, "method", methodName, "host", host, "port", port)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"host":   host,
		"port":   port,
		"method": json.RawMessage(methodJSON),
	})
}
