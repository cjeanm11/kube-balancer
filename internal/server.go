package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	port      int
	allocator *Allocator
	server    *http.Server
	wg        sync.WaitGroup
	testMode  bool
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

func WithAllocator(allocator *Allocator) Option {
	return func(s *Server) {
		s.allocator = allocator
	}
}

func WithTestMode() Option {
	return func(s *Server) {
		s.testMode = true
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on port %d", s.port)
	})

	mux.HandleFunc("/workloads", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var workload Workload
		if err := json.NewDecoder(r.Body).Decode(&workload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := workload.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if s.allocator.testMode {
			if err := s.allocator.TestAllocateResources([]Workload{workload}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if err := s.allocator.SubmitBids([]Workload{workload}); err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Workload submitted successfully"))
	})

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		fmt.Printf("Server starting on port %d...\n", s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()
}

func (s *Server) Stop() {
	if s.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	}

	s.wg.Wait()
	fmt.Println("Server stopped gracefully")
}

