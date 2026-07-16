package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	reqs := rl.requests[ip]
	valid := make([]time.Time, 0)
	for _, t := range reqs {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.requests[ip] = valid
		return false
	}

	rl.requests[ip] = append(valid, now)
	return true
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Auth valida un header X-API-Key contra la variable de entorno API_KEY.
// Si API_KEY no esta configurada, el servicio permite peticiones sin autenticacion
// (util para desarrollo local y pruebas).
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			next(w, r)
			return
		}

		key := r.Header.Get("X-API-Key")
		if key != apiKey {
			writeError(w, http.StatusUnauthorized, "API key invalida")
			return
		}
		next(w, r)
	}
}

func RateLimit(limiter *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = strings.Split(fwd, ",")[0]
			}

			if !limiter.Allow(strings.TrimSpace(ip)) {
				writeError(w, http.StatusTooManyRequests, "Demasiadas peticiones. Intenta de nuevo en un minuto")
				return
			}
			next(w, r)
		}
	}
}
