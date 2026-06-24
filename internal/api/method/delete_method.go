package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/darksuei/kubeRPC/internal/cache"
)

func DeleteMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	deleted, err := cache.Store.HDel("service:"+serviceName, methodName)
	if err != nil {
		slog.Error("delete-method: cache error", "service", serviceName, "method", methodName, "error", err)
		http.Error(w, "Error deleting method", http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		slog.Debug("delete-method: not found", "service", serviceName, "method", methodName)
		http.Error(w, "Service or method not found", http.StatusNotFound)
		return
	}

	slog.Info("method deleted", "service", serviceName, "method", methodName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Method [%s] in service [%s] deleted successfully", methodName, serviceName)
}
