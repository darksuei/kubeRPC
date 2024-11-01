package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

// Structure for storing a method
type Method struct {
	Name        string   `json:"name"`
	Params      []string `json:"params"`
	Description string   `json:"description"`
}

// Structure for storing a service
type RegisterRequest struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	ServiceName string   `json:"service_name"`
	Methods     []Method `json:"methods"`
}

// Redis Client Setup
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", // Update with your Redis address if needed
	Password: "",               // no password set
	DB:       0,                // use default DB
})

func Health(w http.ResponseWriter, r *http.Request) {
	// Check Redis connection
	_, err := rdb.Ping().Result()

	if err != nil {
		log.Printf("Failed to connect to Redis: %s", err)
		http.Error(w, "Redis connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is healthy"))
}

func GetServiceFunction(w http.ResponseWriter, r *http.Request) {
	// Get the service name from the query parameters
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

	// Fetch the service details from Redis using HGetAll
	methodDetails, err := rdb.HGet("service:"+serviceName, methodName).Result()

	if err != nil {
		if err == redis.Nil {
			http.Error(w, "Service or Method not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error retrieving service", http.StatusInternalServerError)
		return
	}

	// Check if the service exists
	if len(methodDetails) == 0 {
		http.Error(w, "Service or Method not found", http.StatusNotFound)
		return
	}

	serviceDetails, err := rdb.HGetAll("service:" + serviceName).Result()

	log.Println(serviceDetails)

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Respond with the service details as JSON
	response, err := json.Marshal(serviceDetails)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// RegisterService Function
// At the moment, this would create a new key in Redis for each service
// and store each method as a field in the hash
// If the service does not already exist, it creates it
// If the service already exists, it will update the existing methods
func RegisterServiceFunction(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println(req.ServiceName)
	log.Println("Registering service:", req)

	redisKey := "service:" + req.ServiceName

	// Set the host
	err := rdb.WithContext(context.Background()).HSet(redisKey, "host", req.Host).Err()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store host in Redis", http.StatusInternalServerError)
		return
	}

	// Set the port
	err = rdb.WithContext(context.Background()).HSet(redisKey, "port", req.Port).Err()
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
		err = rdb.WithContext(context.Background()).HSet(redisKey, method.Name, methodData).Err()

		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to store method in Redis", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service [%s] registered successfully", req.ServiceName)
}

func DeleteServiceFunction(w http.ResponseWriter, r *http.Request) {
	// Get the service name from the query parameters
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

	// Delete the service from Redis
	deleted, err := rdb.HDel("service:"+serviceName, methodName).Result()

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

func main() {
	log.Println("Running kubeRPC API on port 8080..")

	http.HandleFunc("/health", Health)
	http.HandleFunc("/get-service-method", GetServiceFunction)
	http.HandleFunc("/register-service-method", RegisterServiceFunction)
	http.HandleFunc("/delete-service-method", DeleteServiceFunction)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
