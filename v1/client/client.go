// The client implements this to send requests to the server.
package client

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/darksuei/kubeRPC/v1/api"
)

type Client struct {
	Addr string
}

// NewClient creates a new client with the server's Addr
func NewClient(address string) *Client {
	return &Client{Addr: address}
}

// Call connects to the server for each request, sends the request, and receives the response
func (c *Client) Call(method string, params interface{}) (interface{}, error) {
	// Establish a new TCP connection
	conn, err := net.DialTimeout("tcp", c.Addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Prepare the request
	request := api.Request{
		Method: method,
		Params: params,
	}

	// Encode the request and send it over the new connection
	encoder := json.NewEncoder(conn)

	if err := encoder.Encode(&request); err != nil {
		return nil, err
	}

	// Receive and decode the response
	var response api.Response

	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	// Check for errors in the response
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response.Result, nil
}
