package internal

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)


type Allocator struct {
	availableCPU int
	mu           sync.Mutex
}

func NewAllocator(cpu int) *Allocator {
	return &Allocator{availableCPU: cpu}
}


func (a *Allocator) Run() {
	// todo
}

func (a *Allocator) AllocateResources(bids []Workload) {
	a.mu.Lock()
	defer a.mu.Unlock()

	fmt.Println("Processing bids...")
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Priority > bids[j].Priority
	})

	rand.Seed(time.Now().UnixNano())
	for _, workload := range bids {
		if a.availableCPU >= workload.CPU {
			fmt.Printf("Allocating %d vCPUs to %s (Priority %d)\n", workload.CPU, workload.Name, workload.Priority)
			a.availableCPU -= workload.CPU
		} else {
			fmt.Printf("Not enough CPU for %s (Priority %d)\n", workload.Name, workload.Priority)
		}
	}

	fmt.Println("Remaining CPU:", a.availableCPU)
}