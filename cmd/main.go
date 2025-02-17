package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"strconv"
	"kube-balancer/internal"
)

func main() {
    log.Println("Application started")

    cpuFlag := flag.Int("cpu", 10, "Total available CPU")
    bidsFlag := flag.String("bids", "", "Workloads in format name:cpu:memory:priority,name:cpu:memory:priority")
    portFlag := flag.Int("port", 8080, "Server port")

    flag.Parse()

    allocator := internal.NewAllocator(*cpuFlag)
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

            bids = append(bids, internal.Workload{
                Name:     parts[0],
                CPU:      cpu,
                Memory:   memory,
                Priority: priority,
            })
        }

        allocator.AllocateResources(bids)
    }

    server := internal.NewServer(internal.WithPort(*portFlag))
    server.Start()
}