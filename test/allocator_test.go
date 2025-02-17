package test

import (
	"testing"
	"kube-balancer/internal" 
)

func TestAllocator_AllocateResources(t *testing.T) {
	allocator := internal.NewAllocator(100)

	bids := []internal.Workload{
	    {Name: "TaskA", Priority: 2, CPU: 30},
	    {Name: "TaskB", Priority: 1, CPU: 50},
	    {Name: "TaskC", Priority: 3, CPU: 40},
	}	
	allocator.AllocateResources(bids)

	if allocator.GetAvailableCPU() != 30 { 
		t.Errorf("Expected remaining CPU to be 30, got %d", allocator.GetAvailableCPU())
	}
}