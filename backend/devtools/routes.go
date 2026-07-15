package devtools

import (
	"log"
	"os"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes exposes local demo routes that are disabled unless explicitly enabled.
func RegisterRoutes(r chi.Router) {
	if os.Getenv(enableGCPressureEnv) == "true" {
		r.Post("/api/debug/gc-pressure", handleGCPressure)
		log.Println("GC pressure debug route enabled at POST /api/debug/gc-pressure")
	}

	if os.Getenv(enableSchedulerPressureEnv) == "true" {
		r.Post("/api/debug/scheduler-pressure", handleSchedulerPressure)
		log.Println("Scheduler pressure debug route enabled at POST /api/debug/scheduler-pressure")
	}
}
