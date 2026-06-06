package models

import (
    "time"
)

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
    Title       string `json:"title" validate:"required,max=255"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

type UpdateTaskRequest struct {
    Title       string `json:"title" validate:"required,max=255"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

type TaskResponse struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}

func (t *Task) ToResponse() *TaskResponse {
    return &TaskResponse{
        ID:          t.ID,
        Title:       t.Title,
        Description: t.Description,
        Completed:   t.Completed,
        CreatedAt:   t.CreatedAt,
    }
}