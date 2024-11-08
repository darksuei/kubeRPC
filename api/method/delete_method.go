package api

import (
	"fmt"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
)

func DeleteMethod(w http.ResponseWriter, r *http.Request) {
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

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	redisKey := "service:" + serviceName

	deleted, err := config.Rdb.HDel(redisKey, methodName).Result()

	if err != nil {
		http.Error(w, "Error deleting service", http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		http.Error(w, "Service or Method not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Method [%s] in Service [%s] deleted successfully", methodName, serviceName)
}
