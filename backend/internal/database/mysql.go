package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
    "netplay-backend/internal/config"
)

type DB struct {
    *sql.DB
}

func InitDB(cfg *config.Config) (*DB, error) {
    // Create DSN (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
    )
    
    // Open database connection
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Database connected successfully")
    
    return &DB{db}, nil
}

func (db *DB) Close() error {
    return db.DB.Close()
}

// Health check
func (db *DB) Health() error {
    return db.Ping()
}