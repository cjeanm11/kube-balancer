package internal

import (
	"testing"
)

func TestAllocatorSystem(t *testing.T) {
	allocator := NewTestAllocator(100, 2048)

	workloads := []Workload{
		{Name: "SystemTaskA", Priority: 2, CPU: 30, Memory: 512},
		{Name: "SystemTaskB", Priority: 1, CPU: 50, Memory: 1024},
		{Name: "SystemTaskC", Priority: 3, CPU: 40, Memory: 256},
	}

	if err := allocator.TestAllocateResources(workloads); err != nil {
		t.Fatalf("Failed to allocate resources: %v", err)
	}

	if allocator.GetAvailableCPU() != 30 {
		t.Errorf("Expected remaining CPU to be 30, got %d", allocator.GetAvailableCPU())
	}

	if allocator.GetAvailableMemory() != 1280 {
		t.Errorf("Expected remaining memory to be 1280MB, got %d", allocator.GetAvailableMemory())
	}

	metrics := allocator.GetMetrics()
	if metrics.Allocations != 2 {
		t.Errorf("Expected 2 allocations, got %d", metrics.Allocations)
	}
} 