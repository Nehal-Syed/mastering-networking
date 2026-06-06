package repositories

import (
    "database/sql"
    "fmt"
    
    "netplay-backend/internal/database"
    "netplay-backend/internal/models"
)

type TaskRepository struct {
    db *database.DB
}

func NewTaskRepository(db *database.DB) *TaskRepository {
    return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
    query := `
        INSERT INTO tasks (title, description, completed)
        VALUES (?, ?, ?)
    `
    
    result, err := r.db.Exec(query, task.Title, task.Description, task.Completed)
    if err != nil {
        return fmt.Errorf("failed to create task: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }
    
    task.ID = int(id)
    return nil
}

func (r *TaskRepository) GetAll() ([]*models.Task, error) {
    query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM tasks
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query tasks: %w", err)
    }
    defer rows.Close()
    
    var tasks []*models.Task
    for rows.Next() {
        task := &models.Task{}
        err := rows.Scan(
            &task.ID,
            &task.Title,
            &task.Description,
            &task.Completed,
            &task.CreatedAt,
            &task.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan task: %w", err)
        }
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

func (r *TaskRepository) GetByID(id int) (*models.Task, error) {
    query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM tasks
        WHERE id = ?
    `
    
    task := &models.Task{}
    err := r.db.QueryRow(query, id).Scan(
        &task.ID,
        &task.Title,
        &task.Description,
        &task.Completed,
        &task.CreatedAt,
        &task.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("task not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get task: %w", err)
    }
    
    return task, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
    query := `
        UPDATE tasks
        SET title = ?, description = ?, completed = ?
        WHERE id = ?
    `
    
    result, err := r.db.Exec(query, task.Title, task.Description, task.Completed, task.ID)
    if err != nil {
        return fmt.Errorf("failed to update task: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

func (r *TaskRepository) Delete(id int) error {
    query := `DELETE FROM tasks WHERE id = ?`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

func (r *TaskRepository) ToggleComplete(id int, completed bool) error {
    query := `UPDATE tasks SET completed = ? WHERE id = ?`
    
    result, err := r.db.Exec(query, completed, id)
    if err != nil {
        return fmt.Errorf("failed to toggle task: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}