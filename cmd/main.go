package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"strconv"
	"syscall"
	"kube-balancer/internal"
)

func main() {
	log.Println("KubeBalancer started")

	cpuFlag := flag.Int("cpu", 10, "Total available CPU")
	memoryFlag := flag.Int("memory", 1024, "Total available memory in MB")
	bidsFlag := flag.String("bids", "", "Workloads in format name:cpu:memory:priority,name:cpu:memory:priority")
	portFlag := flag.Int("port", 8080, "Server port")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()

	if *debugFlag {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	}

	allocator := internal.NewAllocator(*cpuFlag, *memoryFlag)
	go allocator.Run()

	if *bidsFlag != "" {
		var bids []internal.Workload
		workloads := strings.Split(*bidsFlag, ",")
		
		for _, w := range workloads {
			parts := strings.Split(w, ":")
			if len(parts) != 4 {
				fmt.Println("Invalid workload format. Expected: name:cpu:memory:priority")
				return
			}

			cpu, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error parsing CPU:", err)
				return
			}

			memory, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Error parsing Memory:", err)
				return
			}

			priority, err := strconv.Atoi(parts[3])
			if err != nil {
				fmt.Println("Error parsing Priority:", err)
				return
			}

			workload := internal.Workload{
				Name:     parts[0],
				CPU:      cpu,
				Memory:   memory,
				Priority: priority,
			}

			if err := workload.Validate(); err != nil {
				fmt.Printf("Invalid workload %s: %v\n", workload.Name, err)
				return
			}

			bids = append(bids, workload)
		}

		if err := allocator.AllocateResources(bids); err != nil {
			log.Printf("Error allocating resources: %v", err)
			return
		}
	}

	server := internal.NewServer(
		internal.WithPort(*portFlag),
		internal.WithAllocator(allocator),
	)
	server.Start()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	server.Stop()
	allocator.Stop()
}