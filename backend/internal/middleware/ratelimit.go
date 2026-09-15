package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"opsdesk/internal/utils"
)

type clientRecord struct {
	count       int
	lastSeen    time.Time
	windowStart time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientRecord
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientRecord),
		limit:   limit,
		window:  window,
	}

	// Clean up old entries periodically
	go func() {
		for {
			time.Sleep(window)
			rl.mu.Lock()
			now := time.Now()
			for ip, record := range rl.clients {
				if now.Sub(record.lastSeen) > window*2 {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		rl.mu.Lock()
		now := time.Now()
		record, exists := rl.clients[ip]
		if !exists || now.Sub(record.windowStart) > rl.window {
			rl.clients[ip] = &clientRecord{
				count:       1,
				lastSeen:    now,
				windowStart: now,
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		record.lastSeen = now
		if record.count >= rl.limit {
			rl.mu.Unlock()
			utils.Error(w, http.StatusTooManyRequests, "Too many requests, please try again later")
			return
		}

		record.count++
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
