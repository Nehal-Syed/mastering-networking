package services

import (
    "fmt"
    "log"
    
    "netplay-backend/internal/models"
    "netplay-backend/internal/repositories"
)

type TaskService struct {
    repo *repositories.TaskRepository
}

func NewTaskService(repo *repositories.TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
    // Validation
    if req.Title == "" {
        return nil, fmt.Errorf("title is required")
    }
    
    task := &models.Task{
        Title:       req.Title,
        Description: req.Description,
        Completed:   req.Completed,
    }
    
    if err := s.repo.Create(task); err != nil {
        log.Printf("Error creating task: %v", err)
        return nil, fmt.Errorf("failed to create task")
    }
    
    return task, nil
}

func (s *TaskService) GetAllTasks() ([]*models.Task, error) {
    tasks, err := s.repo.GetAll()
    if err != nil {
        log.Printf("Error getting tasks: %v", err)
        return nil, fmt.Errorf("failed to get tasks")
    }
    
    return tasks, nil
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
    task, err := s.repo.GetByID(id)
    if err != nil {
        if err.Error() == "task not found" {
            return nil, fmt.Errorf("task not found")
        }
        log.Printf("Error getting task by ID: %v", err)
        return nil, fmt.Errorf("failed to get task")
    }
    
    return task, nil
}

func (s *TaskService) UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error) {
    // Check if task exists
    existingTask, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Update fields
    existingTask.Title = req.Title
    existingTask.Description = req.Description
    existingTask.Completed = req.Completed
    
    if err := s.repo.Update(existingTask); err != nil {
        log.Printf("Error updating task: %v", err)
        return nil, fmt.Errorf("failed to update task")
    }
    
    return existingTask, nil
}

func (s *TaskService) DeleteTask(id int) error {
    if err := s.repo.Delete(id); err != nil {
        log.Printf("Error deleting task: %v", err)
        return fmt.Errorf("failed to delete task")
    }
    
    return nil
}

func (s *TaskService) ToggleTaskComplete(id int) (*models.Task, error) {
    task, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    newStatus := !task.Completed
    if err := s.repo.ToggleComplete(id, newStatus); err != nil {
        log.Printf("Error toggling task: %v", err)
        return nil, fmt.Errorf("failed to toggle task status")
    }
    
    task.Completed = newStatus
    return task, nil
}