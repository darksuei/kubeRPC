// The server implements this to recieve and respond to the requests from the client
// cmd/server.go
package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/darksuei/kubeRPC/api"
)

type Server struct {
	Addr     string
	Handlers map[string]func(params interface{}) (interface{}, error)
	mu       sync.RWMutex
}

func NewServer(addr string) *Server {
	return &Server{
		Addr:     addr,
		Handlers: make(map[string]func(params interface{}) (interface{}, error)),
	}
}

func (s *Server) Start() error {
	// listen for connections on the specified Addr
	listener, err := net.Listen("tcp", s.Addr)
	log.Printf("Server listening on %s", s.Addr)

	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn) // Handle connections concurrently
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	var request api.Request
	decoder := json.NewDecoder(conn)

	if err := decoder.Decode(&request); err != nil {
		s.respondWithError(conn, err)
		return
	}

	s.mu.RLock()
	handler, ok := s.Handlers[request.Method]
	s.mu.RUnlock()

	if !ok {
		s.respondWithError(conn, http.ErrAbortHandler)
		return
	}

	result, err := handler(request.Params)

	if err != nil {
		s.respondWithError(conn, err)
		return
	}

	response := api.Response{
		Result: result,
	}
	conn.Write([]byte(jsonResponse(response)))
}

func (s *Server) respondWithError(conn net.Conn, err error) {
	response := api.Response{
		Error: err.Error(),
	}
	responseJSON, _ := json.Marshal(response)
	conn.Write(responseJSON) // Send error response
}

func jsonResponse(response api.Response) string {
	res, _ := json.Marshal(response)
	return string(res)
}

// func (s *Server) Register(method string, handler func(params interface{}) (interface{}, error)) {
//     s.mu.Lock()
//     defer s.mu.Unlock()
//     s.handlers[method] = handler
// }
