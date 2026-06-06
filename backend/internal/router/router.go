package router

import (
    "net/http"
    
    "github.com/gorilla/mux"
    "netplay-backend/internal/config"
    "netplay-backend/internal/database"
    "netplay-backend/internal/handlers"
    "netplay-backend/internal/middleware"
    "netplay-backend/internal/repositories"
    "netplay-backend/internal/services"
)

func SetupRouter(cfg *config.Config, db *database.DB) *mux.Router {
    // Initialize rate limiter
    rateLimiter := middleware.NewRateLimiter(cfg)
    rateLimiter.CleanupVisitors()
    
    // Initialize repository, service, handler
    taskRepo := repositories.NewTaskRepository(db)
    taskService := services.NewTaskService(taskRepo)
    taskHandler := handlers.NewTaskHandler(taskService)
    
    // Create main router
    r := mux.NewRouter()
    
    // Apply global middleware
    r.Use(middleware.Logging())
    r.Use(middleware.CORS(cfg))
    r.Use(middleware.NetworkInfo(cfg))
    r.Use(rateLimiter.RateLimit)
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()
    
    // Task routes
    api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
    api.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
    api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTaskByID).Methods("GET")
    api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
    api.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")
    api.HandleFunc("/tasks/{id:[0-9]+}/toggle", taskHandler.ToggleTaskComplete).Methods("PATCH")
    
    // Utility routes
    api.HandleFunc("/network-info", middleware.NetworkInfoHandler(cfg)).Methods("GET")
    api.HandleFunc("/health", middleware.HealthCheckHandler(db)).Methods("GET")
    
    // Root route
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "NetPlay API is running", "version": "1.0.0"}`))
    }).Methods("GET")
    
    return r
}