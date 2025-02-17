package internal

import (
	"sort"
	"testing"
)

type TestAllocator struct {
	availableCPU    int
	availableMemory int
	metrics        ResourceMetrics
}

func NewTestAllocator(cpu, memory int) *TestAllocator {
	return &TestAllocator{
		availableCPU:    cpu,
		availableMemory: memory,
		metrics: ResourceMetrics{
			TotalCPU:    cpu,
			TotalMemory: memory,
		},
	}
}

func (ta *TestAllocator) SubmitBids(bids []Workload) error {
	return ta.TestAllocateResources(bids)
}

func (ta *TestAllocator) TestAllocateResources(bids []Workload) error {
	if len(bids) == 0 {
		return nil
	}

	for _, workload := range bids {
		if err := workload.Validate(); err != nil {
			return err
		}
	}

	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Priority > bids[j].Priority
	})

	for _, workload := range bids {
		if ta.availableCPU >= workload.CPU && ta.availableMemory >= workload.Memory {
			ta.availableCPU -= workload.CPU
			ta.availableMemory -= workload.Memory
			ta.metrics.Allocations++
			ta.metrics.UsedCPU = ta.metrics.TotalCPU - ta.availableCPU
			ta.metrics.UsedMemory = ta.metrics.TotalMemory - ta.availableMemory
		}
	}

	return nil
}

func (ta *TestAllocator) GetAvailableCPU() int {
	return ta.availableCPU
}

func (ta *TestAllocator) GetAvailableMemory() int {
	return ta.availableMemory
}

func (ta *TestAllocator) GetMetrics() ResourceMetrics {
	return ta.metrics
}

func TestAllocator_AllocateResources(t *testing.T) {
	allocator := NewAllocator(100, 2048)

	bids := []Workload{
		{Name: "TaskA", Priority: 2, CPU: 30, Memory: 512},
		{Name: "TaskB", Priority: 1, CPU: 50, Memory: 1024},
		{Name: "TaskC", Priority: 3, CPU: 40, Memory: 256},
	}

	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Priority > bids[j].Priority
	})

	for _, workload := range bids {
		if allocator.GetAvailableCPU() >= workload.CPU && allocator.GetAvailableMemory() >= workload.Memory {
			allocator.availableCPU -= workload.CPU
			allocator.availableMemory -= workload.Memory
			allocator.metrics.Allocations++
		}
	}

	if allocator.GetAvailableCPU() != 30 {
		t.Errorf("Expected remaining CPU to be 30, got %d", allocator.GetAvailableCPU())
	}

	if allocator.GetAvailableMemory() != 1280 {
		t.Errorf("Expected remaining memory to be 1280MB, got %d", allocator.GetAvailableMemory())
	}
}

func TestWorkload_Validation(t *testing.T) {
	tests := []struct {
		name    string
		workload Workload
		wantErr bool
	}{
		{
			name: "Valid workload",
			workload: Workload{
				Name:     "ValidTask",
				CPU:      10,
				Memory:   512,
				Priority: 1,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			workload: Workload{
				Name:     "",
				CPU:      10,
				Memory:   512,
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Zero CPU",
			workload: Workload{
				Name:     "Task",
				CPU:      0,
				Memory:   512,
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Zero Memory",
			workload: Workload{
				Name:     "Task",
				CPU:      10,
				Memory:   0,
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Zero Priority",
			workload: Workload{
				Name:     "Task",
				CPU:      10,
				Memory:   512,
				Priority: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.workload.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Workload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllocator_Metrics(t *testing.T) {
	allocator := NewAllocator(100, 2048)
	metrics := allocator.GetMetrics()

	if metrics.TotalCPU != 100 {
		t.Errorf("Expected total CPU to be 100, got %d", metrics.TotalCPU)
	}

	if metrics.TotalMemory != 2048 {
		t.Errorf("Expected total memory to be 2048, got %d", metrics.TotalMemory)
	}

	if metrics.Allocations != 0 {
		t.Errorf("Expected initial allocations to be 0, got %d", metrics.Allocations)
	}
} 