package service

import (
	"todo-app/backend/internal/models"
	"todo-app/backend/internal/repository"
)

type TaskService interface {
	GetTasks(userID int64, filter string) ([]models.Task, error)
	AddTask(task *models.Task) (*models.Task, error)
	UpdateTask(task *models.Task) (*models.Task, error)
	ToggleTask(userID, id int64) (*models.Task, error)
	DeleteTask(userID, id int64) error
	ClearCompleted(userID int64) (int64, error)
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) GetTasks(userID int64, filter string) ([]models.Task, error) {
	return s.repo.GetAll(userID, filter)
}

func (s *taskService) AddTask(task *models.Task) (*models.Task, error) {
	return s.repo.Create(task)
}

func (s *taskService) UpdateTask(task *models.Task) (*models.Task, error) {
	return s.repo.Update(task)
}

func (s *taskService) ToggleTask(userID, id int64) (*models.Task, error) {
	t, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	if t.Completed {
		t.Completed = false
		t.CompletedAt = nil
	} else {
		t.Completed = true
		now := (t.CreatedAt)
		t.CompletedAt = &now
	}
	return s.repo.Update(t)
}

func (s *taskService) DeleteTask(userID, id int64) error {
	return s.repo.Delete(userID, id)
}

func (s *taskService) ClearCompleted(userID int64) (int64, error) {
	return s.repo.ClearCompleted(userID)
}
