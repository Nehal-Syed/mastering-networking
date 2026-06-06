package middleware

import (
    "net/http"
    "strings"
    
    "netplay-backend/internal/config"
)

func CORS(cfg *config.Config) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if origin is allowed
            allowed := false
            for _, allowedOrigin := range cfg.AllowedOrigins {
                if strings.TrimSpace(allowedOrigin) == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
            
            w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ","))
            w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ","))
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Max-Age", "3600")
            
            // Handle preflight requests
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}