package internal

import (
	"sort"
	"sync"
	"log"
)

type Allocator struct {
	availableCPU int
	mu           sync.Mutex
	bidQueue     chan []Workload
	quit         chan struct{}
}

func NewAllocator(cpu int) *Allocator {
	return &Allocator{
		availableCPU: cpu,
		bidQueue:     make(chan []Workload),
		quit:         make(chan struct{}),
	}
}

func (a *Allocator) GetAvailableCPU() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    return a.availableCPU
}

func (a *Allocator) Run() {
	log.Println("Allocator started.")
	for {
		select {
		case bids := <-a.bidQueue:
			a.AllocateResources(bids)
		case <-a.quit:
			log.Println("Allocator shutting down.")
			return
		}
	}
}

func (a *Allocator) AllocateResources(bids []Workload) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(bids) == 0 {
		log.Println("No bids to process.")
		return
	}

	log.Println("Processing bids...")
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Priority > bids[j].Priority
	})

	for _, workload := range bids {
		if a.availableCPU >= workload.CPU {
			log.Printf("Allocating %d vCPUs to %s (Priority %d)", workload.CPU, workload.Name, workload.Priority)
			a.availableCPU -= workload.CPU
		} else {
			log.Printf("Insufficient CPU for %s (Priority %d), required: %d, available: %d",
				workload.Name, workload.Priority, workload.CPU, a.availableCPU)
		}
	}

	log.Printf("Allocation complete. Remaining CPU: %d", a.availableCPU)
}

func (a *Allocator) SubmitBids(bids []Workload) {
	a.bidQueue <- bids
}

func (a *Allocator) Stop() {
	close(a.quit)
}