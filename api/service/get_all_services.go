package api

import (
	"encoding/json"
	"net/http"

	"github.com/darksuei/kubeRPC/config"
)

func GetAllServices(w http.ResponseWriter, _ *http.Request) {
	keys, err := config.Rdb.Keys("service:*").Result()
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
		serviceDetails, err := config.Rdb.HGetAll(key).Result()
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
