package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
	"github.com/darksuei/kubeRPC/helpers"
	"github.com/go-redis/redis"
)

func UpdateMethod(w http.ResponseWriter, r *http.Request) {
	var req helpers.Method

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

	methodData, err := json.Marshal(req)

	if err != nil {
		http.Error(w, "Error processing method data", http.StatusInternalServerError)
		return
	}

	if methodName != req.Name {
		deleted, err := config.Rdb.HDel(redisKey, methodName).Result()

		if err != nil {
			http.Error(w, "Error deleting service", http.StatusInternalServerError)
			return
		}

		if deleted == 0 {
			http.Error(w, "Service or Method not found", http.StatusNotFound)
			return
		}
		methodName = req.Name
	}

	err = config.Rdb.WithContext(context.Background()).HSet(redisKey, methodName, methodData).Err()

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store method in Redis", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Method [%s] updated successfully", req.Name)
}
