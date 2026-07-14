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

	"github.com/go-chi/chi/v5"
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

// RegisterRoutes exposes local demo routes that are disabled unless explicitly enabled.
func RegisterRoutes(r chi.Router) {
	if os.Getenv(enableGCPressureEnv) != "true" {
		return
	}

	r.Post("/api/debug/gc-pressure", handleGCPressure)
	log.Println("GC pressure debug route enabled at POST /api/debug/gc-pressure")
}

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

func allocateMemory(sizeMB int) []byte {
	chunk := make([]byte, sizeMB*bytesPerMB)
	for offset := 0; offset < len(chunk); offset += os.Getpagesize() {
		chunk[offset] = byte(offset / bytesPerMB)
	}
	return chunk
}

func releaseMemoryAfter(allocationID uint64, holdDuration time.Duration) {
	time.Sleep(holdDuration)

	gcPressureMu.Lock()
	delete(gcPressureAllocations, allocationID)
	activeLoads := len(gcPressureAllocations)
	gcPressureMu.Unlock()

	runtime.GC()
	log.Printf("released GC pressure allocation: allocation_id=%d active_loads=%d", allocationID, activeLoads)
}
