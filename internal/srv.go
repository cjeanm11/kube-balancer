package internal

import (
	"fmt"
	"net/http"
)

type Server struct {
	port int
}

type Option func(*Server)

func NewServer(options ...Option) *Server {
	server := &Server{}

	for _, option := range options {
		option(server)
	}

	return server
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on port %d", s.port)
	})

	address := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server starting on port %d...\n", s.port)
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

