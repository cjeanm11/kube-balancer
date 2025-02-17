# KubeBalancer

A resource allocation system inspired by economic bidding markets, where workloads compete for CPU and memory resources based on priority-based bidding. Higher priority workloads get preferential access to resources, similar to how higher bids win in an auction.

**Rules**:
- Priority: 1-100 (higher wins)
- Resources: CPU/memory > 0
- Names: Unique required
- Ties: First-come wins

## Run

```bash
# Build and run 
CGO_ENABLED=0 go build ./cmd/main.go
./main -cpu 100 -memory 2048 -port 8080 -debug

# Submit workload
curl -X POST http://localhost:8080/workloads \
  -H "Content-Type: application/json" \
  -d '{"name": "test1", "cpu": 10, "memory": 512, "priority": 50}'
```

## Testing

```bash

CGO_ENABLED=0 go test -v ./...

```

## Status

**Core**: 
  - Priority-based allocation (1-100 scale)
  - Resource validation (CPU/memory > 0)
  - HTTP API endpoints (POST /workloads)
  - Buffered bid queue (100 capacity)
  - Graceful server shutdown

**In Progress**:
  - Error handling (timeout management)
  - Performance (bid processing optimization)
  - Resource tracking (allocation monitoring)
  - Server lifecycle (port conflict handling)

**Future**:
  - Kubernetes integration (CRDs, scheduler)
  - Monitoring (metrics, logging)
  - Scaling (multi-node, auto-scaling)
  - Security (authentication, authorization)
