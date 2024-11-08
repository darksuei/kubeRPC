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

func RegisterMethods(w http.ResponseWriter, r *http.Request) {
	var req helpers.Service

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	redisKey := "service:" + req.ServiceName

	serviceDetails, err := config.Rdb.HGetAll("service:" + req.ServiceName).Result()

	if err != nil {
		if err == redis.Nil {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving service", http.StatusInternalServerError)
		return
	}

	if len(serviceDetails) == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	for _, method := range req.Methods {
		methodData, err := json.Marshal(method)

		if err != nil {
			http.Error(w, "Error processing method data", http.StatusInternalServerError)
			return
		}
		err = config.Rdb.WithContext(context.Background()).HSet(redisKey, method.Name, methodData).Err()

		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to store method in Redis", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Methods registered in service %s successfully", req.ServiceName)
}
