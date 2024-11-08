package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/darksuei/kubeRPC/config"
	"github.com/darksuei/kubeRPC/helpers"
)

func RegisterMethods(w http.ResponseWriter, r *http.Request) {
	log.Print("Registering a service method.")
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

	log.Println(req.ServiceName)
	log.Println("Registering service:", req)

	redisKey := "service:" + req.ServiceName

	// Set the host
	err := config.Rdb.WithContext(context.Background()).HSet(redisKey, "serviceName", req.ServiceName).Err()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store service name in Redis", http.StatusInternalServerError)
		return
	}

	// the host that will be stored wont be local host but rather the kubernetes dns name
	host := req.ServiceName + "." + os.Getenv("NAMESPACE") + ".svc.cluster.local"

	// Set the host
	err = config.Rdb.WithContext(context.Background()).HSet(redisKey, "host", host).Err()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store host in Redis", http.StatusInternalServerError)
		return
	}

	// Set the port
	err = config.Rdb.WithContext(context.Background()).HSet(redisKey, "port", req.Port).Err()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store port in Redis", http.StatusInternalServerError)
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
	fmt.Fprintf(w, "Service [%s] registered successfully", req.ServiceName)
}
