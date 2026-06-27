package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
	"github.com/darksuei/kubeRPC/internal/helpers"
	"github.com/darksuei/kubeRPC/internal/metrics"
)

func RegisterMethods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req helpers.Service
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("register-methods: decode failed", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	slog.Debug("register-methods: request received", "service", req.ServiceName, "count", len(req.Methods))

	if req.ServiceName == "" {
		slog.Error("register-methods: service_name missing in request body")
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}

	serviceDetails, err := cache.Store.HGetAll("service:" + req.ServiceName)
	if err != nil {
		slog.Error("register-methods: cache lookup failed", "service", req.ServiceName, "error", err)
		http.Error(w, "Error retrieving service", http.StatusInternalServerError)
		return
	}

	if len(serviceDetails) == 0 {
		slog.Error("register-methods: service not found: update-service must be called first", "service", req.ServiceName)
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	redisKey := "service:" + req.ServiceName
	for _, method := range req.Methods {
		methodData, err := json.Marshal(method)
		if err != nil {
			slog.Error("register-methods: marshal failed", "method", method.Name, "error", err)
			http.Error(w, "Error processing method data", http.StatusInternalServerError)
			return
		}
		if err := cache.Store.HSet(redisKey, method.Name, methodData); err != nil {
			slog.Error("register-methods: store failed", "method", method.Name, "error", err)
			http.Error(w, "Failed to store method", http.StatusInternalServerError)
			return
		}
		metrics.MethodRegistrations.Inc()
		slog.Info("method registered", "service", req.ServiceName, "method", method.Name)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Methods registered in service %s successfully", req.ServiceName)
}
