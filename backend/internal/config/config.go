package config

import (
    "os"
    "strconv"
    "strings"
    "time"
    
    "github.com/joho/godotenv"
)

type Config struct {
    // Server
    Port        string
    Environment string
    
    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    
    // CORS
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    
    // Rate Limiting
    RateLimitRequests int
    RateLimitDuration time.Duration
    
    // Network
    TrustedProxies []string
}

func LoadConfig() (*Config, error) {
    // Load .env file if it exists
    godotenv.Load()
    
    config := &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("ENVIRONMENT", "development"),
        
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBUser:     getEnv("DB_USER", "netplay_user"),
        DBPassword: getEnv("DB_PASSWORD", "netplay_password"),
        DBName:     getEnv("DB_NAME", "netplay_db"),
        
        AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost,http://localhost:3000"), ","),
        AllowedMethods: strings.Split(getEnv("ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
        AllowedHeaders: strings.Split(getEnv("ALLOWED_HEADERS", "Content-Type,Authorization,X-Requested-With"), ","),
        
        RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
        RateLimitDuration: time.Duration(getEnvInt("RATE_LIMIT_DURATION", 60)) * time.Second,
        
        TrustedProxies: strings.Split(getEnv("TRUSTED_PROXIES", "127.0.0.1,::1"), ","),
    }
    
    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}