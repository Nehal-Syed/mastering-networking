package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"netplay-backend/internal/config"
)

// NetworkInfo middleware to extract and add network information to context
func NetworkInfo(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get real client IP considering proxy headers
			clientIP := getClientIP(r, cfg)

			// Add network info to request context (can be used later)
			r.Header.Set("X-Client-IP", clientIP)

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request, cfg *config.Config) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if isTrustedProxy(ip, cfg) {
				continue
			}
			return ip
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if !isTrustedProxy(xri, cfg) {
			return xri
		}
	}

	// Fallback to remote address
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

func isTrustedProxy(ip string, cfg *config.Config) bool {
	for _, trusted := range cfg.TrustedProxies {
		if strings.TrimSpace(trusted) == ip {
			return true
		}
	}
	return false
}

// NetworkInfoHandler returns network details about the request
func NetworkInfoHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r, cfg)

		info := map[string]interface{}{
			"client_ip":       clientIP,
			"remote_addr":     r.RemoteAddr,
			"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
			"x_real_ip":       r.Header.Get("X-Real-IP"),
			"user_agent":      r.UserAgent(),
			"method":          r.Method,
			"protocol":        r.Proto,
			"host":            r.Host,
			"request_uri":     r.RequestURI,
			"headers":         r.Header,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

// HealthCheckHandler returns service health status
func HealthCheckHandler(db interface{ Health() error }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := map[string]interface{}{
			"status":    "healthy",
			"service":   "netplay-backend",
			"timestamp": time.Now().Unix(),
		}

		// Check database health
		if db != nil {
			if err := db.Health(); err != nil {
				status["status"] = "unhealthy"
				status["database"] = "disconnected"
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				status["database"] = "connected"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
