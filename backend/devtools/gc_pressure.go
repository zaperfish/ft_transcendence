package devtools

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	enableGCPressureEnv = "ENABLE_GC_PRESSURE_TEST"

	defaultGCPressureSizeMB = 128
	maxGCPressureSizeMB     = 256
	defaultGCHoldSeconds    = 30
	maxGCHoldSeconds        = 120

	bytesPerMB = 1024 * 1024
)

var (
	gcPressureID          atomic.Uint64
	gcPressureAllocations = make(map[uint64]gcPressureAllocation)
	gcPressureMu          sync.Mutex
)

type gcPressureAllocation struct {
	bytes []byte
}

type gcPressureResponse struct {
	AllocationID uint64 `json:"allocation_id"`
	AllocatedMB  int    `json:"allocated_mb,omitempty"`
	HoldSeconds  int    `json:"hold_seconds"`
	ActiveLoads  int    `json:"active_loads"`
}

// handleGCPressure starts a temporary heap allocation from query parameters.
func handleGCPressure(w http.ResponseWriter, r *http.Request) {
	sizeMB := boundedQueryInt(r, "size_mb", defaultGCPressureSizeMB, 1, maxGCPressureSizeMB)
	holdSeconds := boundedQueryInt(r, "hold_seconds", defaultGCHoldSeconds, 1, maxGCHoldSeconds)

	allocation := gcPressureAllocation{
		bytes: allocateMemory(sizeMB),
	}
	response := gcPressureResponse{
		AllocatedMB: sizeMB,
		HoldSeconds: holdSeconds,
	}

	allocationID := gcPressureID.Add(1)

	gcPressureMu.Lock()
	gcPressureAllocations[allocationID] = allocation
	activeLoads := len(gcPressureAllocations)
	gcPressureMu.Unlock()

	go releaseMemoryAfter(allocationID, time.Duration(holdSeconds)*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response.AllocationID = allocationID
	response.ActiveLoads = activeLoads
	_ = json.NewEncoder(w).Encode(response)
}

// boundedQueryInt reads an integer query parameter and clamps it within the allowed range.
func boundedQueryInt(r *http.Request, name string, defaultValue int, minValue int, maxValue int) int {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

// allocateMemory allocates and touches a byte slice so it is visible to the runtime and OS.
func allocateMemory(sizeMB int) []byte {
	chunk := make([]byte, sizeMB*bytesPerMB)
	for offset := 0; offset < len(chunk); offset += os.Getpagesize() {
		chunk[offset] = byte(offset / bytesPerMB)
	}
	return chunk
}

// releaseMemoryAfter removes a tracked allocation after its hold duration and forces a GC cycle.
func releaseMemoryAfter(allocationID uint64, holdDuration time.Duration) {
	time.Sleep(holdDuration)

	gcPressureMu.Lock()
	delete(gcPressureAllocations, allocationID)
	activeLoads := len(gcPressureAllocations)
	gcPressureMu.Unlock()

	runtime.GC()
	log.Printf("released GC pressure allocation: allocation_id=%d active_loads=%d", allocationID, activeLoads)
}
