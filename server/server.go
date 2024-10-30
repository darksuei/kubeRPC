// The server implements this to recieve and respond to the requests from the client
// cmd/server.go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/darksueii/kubeRPC/pkg/api"
)

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	var req api.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp api.Response
	switch req.Method {
	case "exampleMethod":
		// Example response processing
		resp.Result = "Hello from kubeRPC!"
	default:
		resp.Error = "Method not found"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/rpc", rpcHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
