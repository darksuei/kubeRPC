package api

import (
	"encoding/json"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
	"github.com/go-redis/redis"
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

	methodDetails, err := config.Rdb.HGet("service:"+serviceName, methodName).Result()

	if err != nil {
		if err == redis.Nil {
			http.Error(w, "Service or Method not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving service", http.StatusInternalServerError)
		return
	}

	if len(methodDetails) == 0 {
		http.Error(w, "Service or Method not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	serviceDetails, err := config.Rdb.HGetAll("service:" + serviceName).Result()

	response, err := json.Marshal(serviceDetails)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
