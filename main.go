package main

import (
	"log"
	"net/http"

	api "github.com/darksuei/kubeRPC/api"
	method "github.com/darksuei/kubeRPC/api/method"
	service "github.com/darksuei/kubeRPC/api/service"
	config "github.com/darksuei/kubeRPC/config"
	serviceDiscovery "github.com/darksuei/kubeRPC/service_discovery"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.InitRedisClient()

	clientset, err := serviceDiscovery.CreateKubeClient()

	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	serviceDiscovery.GetKubeServices(clientset)

	log.Println("Running kubeRPC API on port 8080..")

	http.HandleFunc("/health", api.Health)
	http.HandleFunc("/get-service", service.GetService)
	http.HandleFunc("/get-all-services", service.GetServices)
	http.HandleFunc("/get-service-method", method.GetMethod)
	http.HandleFunc("/register-service-method", method.RegisterMethods)
	http.HandleFunc("/delete-service-method", method.DeleteMethod)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
