package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
	"github.com/darksuei/kubeRPC/internal/helpers"
)

func UpdateService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	var req helpers.Service
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("update-service: decode failed", "service", serviceName, "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	slog.Debug("update-service: registering", "service", serviceName, "host", req.Host, "port", req.Port)

	redisKey := "service:" + serviceName

	if err := cache.Store.HSet(redisKey, "serviceName", serviceName); err != nil {
		slog.Error("update-service: store serviceName failed", "service", serviceName, "error", err)
		http.Error(w, "Failed to store service", http.StatusInternalServerError)
		return
	}

	if req.Host != "" {
		if err := cache.Store.HSet(redisKey, "host", req.Host); err != nil {
			slog.Error("update-service: store host failed", "service", serviceName, "error", err)
			http.Error(w, "Failed to store host", http.StatusInternalServerError)
			return
		}
	}

	if req.Port != 0 {
		if err := cache.Store.HSet(redisKey, "port", req.Port); err != nil {
			slog.Error("update-service: store port failed", "service", serviceName, "error", err)
			http.Error(w, "Failed to store port", http.StatusInternalServerError)
			return
		}
	}

	slog.Info("service registered", "service", serviceName, "host", req.Host, "port", req.Port)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service %s updated successfully", serviceName)
}
