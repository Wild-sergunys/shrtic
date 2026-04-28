package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LoginAttempt хранит информацию о попытках входа с одного IP
type LoginAttempt struct {
	Count        int
	FirstTry     time.Time
	BlockedUntil time.Time
}

// RateLimiter отслеживает попытки входа
type RateLimiter struct {
	mu          sync.RWMutex
	attempts    map[string]*LoginAttempt
	maxAttempts int
	window      time.Duration
	blockTime   time.Duration
}

func NewRateLimiter(maxAttempts int, window, blockTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts:    make(map[string]*LoginAttempt),
		maxAttempts: maxAttempts,
		window:      window,
		blockTime:   blockTime,
	}

	go rl.cleanup()
	return rl
}

// cleanup периодически удаляет устаревшие записи
func (rl *RateLimiter) cleanup() {
	interval := rl.window
	if interval < time.Minute {
		interval = time.Minute
	}
	if interval > time.Hour {
		interval = time.Hour
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, attempt := range rl.attempts {
			if !attempt.BlockedUntil.IsZero() && now.After(attempt.BlockedUntil) {
				if now.Sub(attempt.BlockedUntil) > rl.blockTime {
					delete(rl.attempts, ip)
					continue
				}
			}
			if now.Sub(attempt.FirstTry) > rl.window+rl.blockTime {
				delete(rl.attempts, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow проверяет, можно ли попытаться войти с этого IP
func (rl *RateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	attempt, exists := rl.attempts[ip]
	if !exists {
		rl.attempts[ip] = &LoginAttempt{
			Count:    1,
			FirstTry: now,
		}
		return true, 0
	}

	if !attempt.BlockedUntil.IsZero() {
		if now.Before(attempt.BlockedUntil) {
			remaining := attempt.BlockedUntil.Sub(now)
			return false, remaining
		}
		attempt.BlockedUntil = time.Time{}
		attempt.Count = 1
		attempt.FirstTry = now
		return true, 0
	}

	if now.Sub(attempt.FirstTry) > rl.window {
		attempt.Count = 1
		attempt.FirstTry = now
		return true, 0
	}

	attempt.Count++

	if attempt.Count > rl.maxAttempts {
		attempt.BlockedUntil = now.Add(rl.blockTime)
		return false, rl.blockTime
	}

	return true, 0
}

// Reset сбрасывает счётчик попыток для IP
func (rl *RateLimiter) Reset(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.attempts, ip)
}

// GetIP получает IP клиента из запроса
func GetIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// LoginRateLimitMiddleware ограничивает попытки входа
func LoginRateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetIP(r)

			allowed, remaining := rl.Allow(ip)
			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "rate_limit",
					"message": "Слишком много попыток входа. Попробуйте позже.",
					"details": map[string]interface{}{
						"blocked_for_seconds": int(remaining.Seconds()),
					},
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
