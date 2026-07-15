package devtools

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	enableSchedulerPressureEnv = "ENABLE_SCHEDULER_PRESSURE_TEST"

	schedulerPressureModeSleep = "sleep"
	schedulerPressureModeCPU   = "cpu"

	defaultSchedulerGoroutines  = 1000
	maxSchedulerGoroutines      = 5000
	defaultSchedulerHoldSeconds = 30
	maxSchedulerHoldSeconds     = 120
)

var (
	schedulerPressureID    atomic.Uint64
	schedulerPressureLoads = make(map[uint64]schedulerPressureLoad)
	schedulerPressureMu    sync.Mutex
	schedulerPressureSink  atomic.Uint64
)

type schedulerPressureLoad struct {
	done       chan struct{}
	mode       string
	goroutines int
}

type schedulerPressureResponse struct {
	LoadID           uint64 `json:"load_id"`
	Mode             string `json:"mode"`
	Goroutines       int    `json:"goroutines"`
	HoldSeconds      int    `json:"hold_seconds"`
	ActiveLoads      int    `json:"active_loads"`
	ActiveGoroutines int    `json:"active_goroutines"`
}

// handleSchedulerPressure starts a temporary scheduler pressure load from query parameters.
func handleSchedulerPressure(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	switch mode {
	case "":
		mode = schedulerPressureModeSleep
	case schedulerPressureModeSleep, schedulerPressureModeCPU:
	default:
		http.Error(w, "mode must be one of: sleep, cpu", http.StatusBadRequest)
		return
	}

	goroutineCount := boundedQueryInt(r, "goroutines", defaultSchedulerGoroutines, 1, maxSchedulerGoroutines)
	holdSeconds := boundedQueryInt(r, "hold_seconds", defaultSchedulerHoldSeconds, 1, maxSchedulerHoldSeconds)

	loadID := schedulerPressureID.Add(1)
	load := schedulerPressureLoad{
		done:       make(chan struct{}),
		mode:       mode,
		goroutines: goroutineCount,
	}

	schedulerPressureMu.Lock()
	schedulerPressureLoads[loadID] = load
	activeLoads := len(schedulerPressureLoads)
	activeGoroutines := activeSchedulerPressureGoroutinesLocked()
	schedulerPressureMu.Unlock()

	for i := 0; i < goroutineCount; i++ {
		go runSchedulerPressureWorker(mode, load.done)
	}
	go releaseSchedulerPressureAfter(loadID, time.Duration(holdSeconds)*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(schedulerPressureResponse{
		LoadID:           loadID,
		Mode:             mode,
		Goroutines:       goroutineCount,
		HoldSeconds:      holdSeconds,
		ActiveLoads:      activeLoads,
		ActiveGoroutines: activeGoroutines,
	})
}

// runSchedulerPressureWorker keeps one pressure goroutine alive until the load is released.
func runSchedulerPressureWorker(mode string, done <-chan struct{}) {
	if mode == schedulerPressureModeSleep {
		<-done
		return
	}

	var value uint64
	for {
		select {
		case <-done:
			schedulerPressureSink.Add(value)
			return
		default:
		}

		for i := uint64(0); i < 10000; i++ {
			value += i
			value ^= value << 7
		}
	}
}

// releaseSchedulerPressureAfter stops a pressure load after its configured hold duration.
func releaseSchedulerPressureAfter(loadID uint64, holdDuration time.Duration) {
	time.Sleep(holdDuration)

	schedulerPressureMu.Lock()
	load, ok := schedulerPressureLoads[loadID]
	if ok {
		delete(schedulerPressureLoads, loadID)
	}
	activeLoads := len(schedulerPressureLoads)
	activeGoroutines := activeSchedulerPressureGoroutinesLocked()
	schedulerPressureMu.Unlock()

	if !ok {
		return
	}

	close(load.done)
	log.Printf(
		"released scheduler pressure load: load_id=%d mode=%s goroutines=%d active_loads=%d active_goroutines=%d",
		loadID,
		load.mode,
		load.goroutines,
		activeLoads,
		activeGoroutines,
	)
}

// activeSchedulerPressureGoroutinesLocked counts goroutines tracked in active scheduler pressure loads.
func activeSchedulerPressureGoroutinesLocked() int {
	activeGoroutines := 0
	for _, load := range schedulerPressureLoads {
		activeGoroutines += load.goroutines
	}
	return activeGoroutines
}
