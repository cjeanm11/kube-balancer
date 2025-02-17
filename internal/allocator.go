package internal

import (
	"sort"
	"sync"
	"log"
	"time"
	"fmt"
)

type ResourceMetrics struct {
	TotalCPU     int
	TotalMemory  int
	UsedCPU      int
	UsedMemory   int
	Allocations  int
	LastUpdated  time.Time
}

type Allocator struct {
	availableCPU    int
	availableMemory int
	mu             sync.Mutex
	bidQueue       chan []Workload
	quit           chan struct{}
	done           chan struct{}
	metrics        ResourceMetrics
	testMode       bool
}

func NewAllocator(cpu, memory int) *Allocator {
	return &Allocator{
		availableCPU:    cpu,
		availableMemory: memory,
		bidQueue:        make(chan []Workload, 100),
		quit:            make(chan struct{}),
		done:            make(chan struct{}),
		metrics: ResourceMetrics{
			TotalCPU:    cpu,
			TotalMemory: memory,
			LastUpdated: time.Now(),
		},
	}
}

func (a *Allocator) SetTestMode(enabled bool) {
	a.testMode = enabled
}

func (a *Allocator) GetAvailableCPU() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.availableCPU
}

func (a *Allocator) GetAvailableMemory() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.availableMemory
}

func (a *Allocator) GetMetrics() ResourceMetrics {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.metrics
}

func (a *Allocator) Run() {
	log.Println("Allocator started.")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case bids := <-a.bidQueue:
			if a.testMode {
				a.TestAllocateResources(bids)
			} else {
				a.AllocateResources(bids)
			}
			select {
			case <-a.done:
				// Channel already closed
			default:
				close(a.done)
			}
			a.done = make(chan struct{})
		case <-ticker.C:
			a.updateMetrics()
		case <-a.quit:
			log.Println("Allocator shutting down.")
			close(a.done)
			return
		}
	}
}

func (a *Allocator) updateMetrics() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.metrics.UsedCPU = a.metrics.TotalCPU - a.availableCPU
	a.metrics.UsedMemory = a.metrics.TotalMemory - a.availableMemory
	a.metrics.LastUpdated = time.Now()
}

func (a *Allocator) AllocateResources(bids []Workload) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(bids) == 0 {
		log.Println("No bids to process.")
		return nil
	}

	for _, workload := range bids {
		if err := workload.Validate(); err != nil {
			return err
		}
	}

	log.Println("Processing bids...")
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Priority > bids[j].Priority
	})

	for _, workload := range bids {
		if a.availableCPU >= workload.CPU && a.availableMemory >= workload.Memory {
			log.Printf("Allocating %d vCPUs and %dMB memory to %s (Priority %d)", 
				workload.CPU, workload.Memory, workload.Name, workload.Priority)
			a.availableCPU -= workload.CPU
			a.availableMemory -= workload.Memory
			a.metrics.Allocations++
		} else {
			log.Printf("Insufficient resources for %s (Priority %d). Required: %d vCPUs, %dMB memory. Available: %d vCPUs, %dMB memory",
				workload.Name, workload.Priority, workload.CPU, workload.Memory, 
				a.availableCPU, a.availableMemory)
		}
	}

	a.updateMetrics()
	log.Printf("Allocation complete. Remaining: %d vCPUs, %dMB memory", 
		a.availableCPU, a.availableMemory)
	return nil
}

func (a *Allocator) SubmitBids(bids []Workload) error {
	select {
	case a.bidQueue <- bids:
		return nil
	case <-time.After(5 * time.Second):
		log.Println("Bid submission timed out")
		return fmt.Errorf("bid submission timed out")
	}
}

func (a *Allocator) Wait() {
	<-a.done
}

func (a *Allocator) Stop() {
	close(a.quit)
	close(a.bidQueue)
	select {
	case <-a.done:
	default:
		close(a.done)
	}
}

func (a *Allocator) TestAllocateResources(bids []Workload) error {
	return a.AllocateResources(bids)
}