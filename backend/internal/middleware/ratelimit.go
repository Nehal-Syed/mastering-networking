package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "netplay-backend/internal/config"
    "golang.org/x/time/rate"
)

type RateLimiter struct {
    visitors map[string]*rate.Limiter
    mu       sync.RWMutex
    requests int
    duration time.Duration
}

func NewRateLimiter(cfg *config.Config) *RateLimiter {
    return &RateLimiter{
        visitors: make(map[string]*rate.Limiter),
        requests: cfg.RateLimitRequests,
        duration: cfg.RateLimitDuration,
    }
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.visitors[ip]
    if !exists {
        // Allow requests per duration
        limiter = rate.NewLimiter(rate.Every(rl.duration/time.Duration(rl.requests)), rl.requests)
        rl.visitors[ip] = limiter
    }
    
    return limiter
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        
        limiter := rl.getLimiter(ip)
        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Cleanup old visitors periodically
func (rl *RateLimiter) CleanupVisitors() {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            rl.mu.Lock()
            rl.visitors = make(map[string]*rate.Limiter)
            rl.mu.Unlock()
        }
    }()
}