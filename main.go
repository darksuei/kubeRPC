package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// Channel to signal graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Initialize Kubernetes client and watchers
	if os.Getenv("ENABLE_DEFAULT_SERVICE_DISCOVERY") == "true" {
		clientset, err := serviceDiscovery.CreateKubeClient()
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %v", err)
		}

		go serviceDiscovery.GetKubeServices(clientset)
		go serviceDiscovery.WatchKubeServices(clientset)
	}

	// Health API
	http.HandleFunc("/health", api.Health)

	// Service API
	http.HandleFunc("/get-service", service.GetService)
	http.HandleFunc("/get-services", service.GetServices)
	http.HandleFunc("/update-service", service.UpdateService)
	http.HandleFunc("/delete-service", service.DeleteService)

	// Method API
	http.HandleFunc("/register-methods", method.RegisterMethods)
	http.HandleFunc("/get-method", method.GetMethod)
	http.HandleFunc("/get-methods", method.GetMethods)
	http.HandleFunc("/update-method", method.UpdateMethod)
	http.HandleFunc("/delete-method", method.DeleteMethod)

	// Start HTTP server in a goroutine
	server := &http.Server{Addr: ":8080"}
	go func() {
		log.Println("Running kubeRPC API on port 8080!")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("Shutting down server...")

	// Gracefully shutdown HTTP server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %s", err)
	}

	log.Println("Server exited")
}
