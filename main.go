package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/darksuei/kubeRPC/service_discovery"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
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
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func Health(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check!")
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

func getServiceMethod(w http.ResponseWriter, r *http.Request) {
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

	methodDetails, err := rdb.HGet("service:"+serviceName, methodName).Result()

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

	// Still return the service details because the host and port are needed
	serviceDetails, err := rdb.HGetAll("service:" + serviceName).Result()

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(serviceDetails)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func registerServiceMethods(w http.ResponseWriter, r *http.Request) {
	log.Print("Registering a service method.")
	var req RegisterRequest

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
	err := rdb.WithContext(context.Background()).HSet(redisKey, "serviceName", req.ServiceName).Err()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to store service name in Redis", http.StatusInternalServerError)
		return
	}

	// the host that will be stored wont be local host but rather the kubernetes dns name
	host := req.ServiceName + "." + os.Getenv("NAMESPACE") + ".svc.cluster.local"

	log.Print(host)

	// Set the host
	err = rdb.WithContext(context.Background()).HSet(redisKey, "host", host).Err()
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

func deleteServiceMethod(w http.ResponseWriter, r *http.Request) {
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

func getService(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	serviceDetails, err := rdb.HGetAll("service:" + serviceName).Result()

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

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(serviceDetails)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func getAllServices(w http.ResponseWriter, _ *http.Request) {
	keys, err := rdb.Keys("service:*").Result()
	if err != nil {
		http.Error(w, "Error retrieving services", http.StatusInternalServerError)
		return
	}

	if len(keys) == 0 {
		http.Error(w, "No services found", http.StatusNotFound)
		return
	}

	services := make(map[string]interface{})

	for _, key := range keys {
		serviceDetails, err := rdb.HGetAll(key).Result()
		if err != nil {
			http.Error(w, "Error retrieving service", http.StatusInternalServerError)
			return
		}

		services[key] = serviceDetails
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(services)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	godotenv.Load()

	clientset, err := service_discovery.CreateKubeClient()

	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	service_discovery.GetKubeServices(clientset)

	log.Println("Running kubeRPC API on port 8080..")

	http.HandleFunc("/health", Health)
	http.HandleFunc("/get-service", getService)
	http.HandleFunc("/get-all-services", getAllServices)
	http.HandleFunc("/get-service-method", getServiceMethod)
	http.HandleFunc("/register-service-method", registerServiceMethods)
	http.HandleFunc("/delete-service-method", deleteServiceMethod)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
