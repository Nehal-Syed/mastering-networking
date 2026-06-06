-- Create database
CREATE DATABASE IF NOT EXISTS netplay_db;
USE netplay_db;

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_completed (completed),
    INDEX idx_created_at (created_at)
);


-- Show tables
SHOW TABLES;
SELECT * FROM tasks;