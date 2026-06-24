package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
	"github.com/darksuei/kubeRPC/internal/helpers"
)

func UpdateMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	var req helpers.Method
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("update-method: decode failed", "service", serviceName, "method", methodName, "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	redisKey := "service:" + serviceName

	if _, err := cache.Store.HGet(redisKey, methodName); err != nil {
		if err == cache.ErrNotFound {
			slog.Debug("update-method: not found", "service", serviceName, "method", methodName)
			http.Error(w, "Service or method not found", http.StatusNotFound)
			return
		}
		slog.Error("update-method: cache error", "service", serviceName, "method", methodName, "error", err)
		http.Error(w, "Error retrieving method", http.StatusInternalServerError)
		return
	}

	if methodName != req.Name {
		if deleted, err := cache.Store.HDel(redisKey, methodName); err != nil || deleted == 0 {
			slog.Error("update-method: rename failed", "service", serviceName, "old", methodName, "new", req.Name)
			http.Error(w, "Error updating method", http.StatusInternalServerError)
			return
		}
		slog.Debug("update-method: renamed", "service", serviceName, "old", methodName, "new", req.Name)
		methodName = req.Name
	}

	methodData, err := json.Marshal(req)
	if err != nil {
		slog.Error("update-method: marshal failed", "service", serviceName, "method", methodName, "error", err)
		http.Error(w, "Error processing method data", http.StatusInternalServerError)
		return
	}

	if err := cache.Store.HSet(redisKey, methodName, methodData); err != nil {
		slog.Error("update-method: store failed", "service", serviceName, "method", methodName, "error", err)
		http.Error(w, "Failed to store method", http.StatusInternalServerError)
		return
	}

	slog.Info("method updated", "service", serviceName, "method", req.Name)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Method [%s] updated successfully", req.Name)
}
