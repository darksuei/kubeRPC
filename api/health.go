package api

import (
	"log"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := config.Rdb.Ping().Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %s", err)
		http.Error(w, "Redis connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("kubeRPC core is healthy!"))
}
