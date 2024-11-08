package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
	"github.com/darksuei/kubeRPC/helpers"
)

// Given a service name, update the service details
// If the service does not exist, create it.

func UpdateService(w http.ResponseWriter, r *http.Request) {
	var req helpers.Service

	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	redisKey := "service:" + serviceName

	// Set the service name
	err := config.Rdb.WithContext(context.Background()).HSet(redisKey, "serviceName", serviceName).Err()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store service name in Redis", http.StatusInternalServerError)
		return
	}

	if req.Host != "" {
		// Set the service host
		err = config.Rdb.WithContext(context.Background()).HSet(redisKey, "host", req.Host).Err()
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to store host in Redis", http.StatusInternalServerError)
			return
		}
	}

	if req.Port != 0 {
		// Set the service port
		err = config.Rdb.WithContext(context.Background()).HSet(redisKey, "port", req.Port).Err()
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to store port in Redis", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service %s updated successfully", serviceName)
}
