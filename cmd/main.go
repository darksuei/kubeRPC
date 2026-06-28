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
	"github.com/darksuei/kubeRPC/internal/webhook"
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

	tlsCert := os.Getenv("TLS_CERT_PATH")
	if tlsCert == "" {
		tlsCert = "/etc/kuberpc/tls/tls.crt"
	}

	tlsKey := os.Getenv("TLS_KEY_PATH")
	if tlsKey == "" {
		tlsKey = "/etc/kuberpc/tls/tls.key"
	}

	var webhookSrv *http.Server

	if certsExist(tlsCert, tlsKey) {
		webhookPort := os.Getenv("WEBHOOK_PORT")
		if webhookPort == "" {
			webhookPort = "9443"
		}

		webhookMux := http.NewServeMux()
		webhookMux.HandleFunc("/mutate", webhook.Mutate)

		webhookSrv = &http.Server{
			Addr:    ":" + webhookPort,
			Handler: api.LoggingMiddleware(webhookMux),
		}

		go func() {
			slog.Info("webhook server listening", "port", webhookPort)
			if err := webhookSrv.ListenAndServeTLS(tlsCert, tlsKey); err != nil && err != http.ErrServerClosed {
				slog.Error("webhook server failed", "error", err)
				os.Exit(1)
			}
		}()
	} else {
		slog.Info("webhook server disabled (TLS certs not found)", "cert", tlsCert, "key", tlsKey)
	}

	<-stop
	slog.Info("shutting down")

	ctx := context.Background()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("api server shutdown failed", "error", err)
	}
	if webhookSrv != nil {
		if err := webhookSrv.Shutdown(ctx); err != nil {
			slog.Error("webhook server shutdown failed", "error", err)
		}
	}

	slog.Info("server exited")
}

func certsExist(certPath, keyPath string) bool {
	for _, p := range []string{certPath, keyPath} {
		if _, err := os.Stat(p); err != nil {
			return false
		}
	}
	return true
}
