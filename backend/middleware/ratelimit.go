package middleware

import (
	// Std
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	// External
	"golang.org/x/time/rate"
)

type LimiterStore struct {
	mu sync.Mutex

	// These maps can grow infinitely over the use of the program
	// Better would be to use something like a LRU cache
	IpLimiters   map[string]*rate.Limiter
	UserLimiters map[string]*rate.Limiter
}

func getClientIp(r *http.Request) string {
	// First try to get client ip address
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	// Fallback
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Last fallback (this doesnt really do anything tho)
	// Should be the ip of the reverse proxy
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// I dunno what to do here
	}

	return host
}

func RateLimiterMiddleware(store *LimiterStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := getClientIp(r)
			fmt.Printf("RATE LIMITER IP = %s\n", ip)

			store.mu.Lock()

			// User rate limiting is not implemented yet since the context doesnt
			// provide the user ID yet

			// limiter, exists := store.UserLimiters[userID]
			// if !exists {
			// 	store.UserLimiters[userID] = rate.NewLimiter(5, 20)
			// }

			limiter, exists := store.IpLimiters[ip]
			if !exists {
				limiter = rate.NewLimiter(5, 20)
				store.IpLimiters[ip] = limiter
			}

			if !limiter.Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				store.mu.Unlock()
				return
			}

			store.mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
