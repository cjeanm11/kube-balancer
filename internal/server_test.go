package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerWorkloadSubmission(t *testing.T) {
	allocator := NewTestAllocator(100, 2048)

	workload := Workload{
		Name:     "test-workload",
		CPU:      1,
		Memory:   1024,
		Priority: 1,
	}

	payload, err := json.Marshal(workload)
	if err != nil {
		t.Fatalf("Failed to marshal workload: %v", err)
	}

	req := httptest.NewRequest("POST", "/workloads", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mux := http.NewServeMux()
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

		if err := allocator.TestAllocateResources([]Workload{workload}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	metrics := allocator.GetMetrics()
	if metrics.Allocations != 1 {
		t.Errorf("Expected 1 allocation, got %d", metrics.Allocations)
	}
} 