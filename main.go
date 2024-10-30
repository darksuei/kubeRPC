package main

import (
	"fmt"
	"log"
	"time"

	"github.com/darksuei/kubeRPC/client"
	"github.com/darksuei/kubeRPC/server"
)

// Basically
func main() {
	server := server.Server{
		Addr: ":8080", // Set the address to listen on
		Handlers: map[string]func(interface{}) (interface{}, error){
			"testMethod": testHandler,
		},
	}

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait a moment for the server to start
	time.Sleep(3 * time.Second)

	// Set up the client
	client := client.NewClient("localhost:8080")

	// Test client call to the server
	params := []float64{3, 5}
	
	result, err := client.Call("testMethod", params)
	if err != nil {
		log.Fatalf("Client call error: %v", err)
	}

	// Display the result
	fmt.Printf("Result of addition: %v\n", result)
}

// Test handler function
func testHandler(params interface{}) (interface{}, error) {
	return fmt.Sprintf("Received: %v", params), nil
}
