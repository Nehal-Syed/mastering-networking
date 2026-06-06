package main

import (
    "log"
    "net/http"
    "os"
    
    "netplay-backend/internal/config"
    "netplay-backend/internal/database"
    "netplay-backend/internal/router"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize database
    db, err := database.InitDB(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()
    
    // Setup router
    r := router.SetupRouter(cfg, db)
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on :%s", port)
    log.Printf("Environment: %s", cfg.Environment)
    
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}