package middleware

import (
    "log"
    "net/http"
    "time"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func Logging() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer to capture status code
            wrapped := &responseWriter{
                ResponseWriter: w,
                statusCode:    http.StatusOK,
            }
            
            // Log request details
            log.Printf("[REQUEST] %s %s %s from %s", 
                r.Method, 
                r.URL.Path, 
                r.Proto, 
                r.RemoteAddr)
            
            // Process request
            next.ServeHTTP(wrapped, r)
            
            // Log response details
            duration := time.Since(start)
            log.Printf("[RESPONSE] %s %s - Status: %d - Duration: %v - User-Agent: %s",
                r.Method,
                r.URL.Path,
                wrapped.statusCode,
                duration,
                r.UserAgent())
        })
    }
}

// Detailed logging for network debugging
func NetworkLogging() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Printf("=== Network Debug ===")
            log.Printf("Remote Address: %s", r.RemoteAddr)
            log.Printf("Forwarded For: %s", r.Header.Get("X-Forwarded-For"))
            log.Printf("Real IP: %s", r.Header.Get("X-Real-IP"))
            log.Printf("Request ID: %s", r.Header.Get("X-Request-ID"))
            log.Printf("Headers: %v", r.Header)
            
            next.ServeHTTP(w, r)
        })
    }
}