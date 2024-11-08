package api

import (
	"fmt"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
)

func DeleteService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deleted, err := config.Rdb.Del("service:" + serviceName).Result()

	if err != nil {
		http.Error(w, "Error deleting service", http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service %s deleted successfully", serviceName)
}
