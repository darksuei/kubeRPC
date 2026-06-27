package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "github.com/darksuei/kubeRPC/internal/api"
	method "github.com/darksuei/kubeRPC/internal/api/method"
	service "github.com/darksuei/kubeRPC/internal/api/service"
	"github.com/darksuei/kubeRPC/internal/cache"
	"github.com/darksuei/kubeRPC/internal/logger"
	_ "github.com/darksuei/kubeRPC/internal/metrics"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	godotenv.Load()
	logger.Init()

	slog.Info("kubeRPC core starting",
		"log_level", os.Getenv("LOG_LEVEL"),
		"cache_type", os.Getenv("CACHE_TYPE"),
	)

	cache.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", api.Health)
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/get-service", service.GetService)
	mux.HandleFunc("/get-services", service.GetServices)
	mux.HandleFunc("/update-service", service.UpdateService)
	mux.HandleFunc("/delete-service", service.DeleteService)

	mux.HandleFunc("/register-methods", method.RegisterMethods)
	mux.HandleFunc("/get-method", method.GetMethod)
	mux.HandleFunc("/get-methods", method.GetMethods)
	mux.HandleFunc("/update-method", method.UpdateMethod)
	mux.HandleFunc("/delete-method", method.DeleteMethod)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: api.LoggingMiddleware(mux),
	}

	go func() {
		slog.Info("server listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		slog.Error("shutdown failed", "error", err)
	}
	slog.Info("server exited")
}
