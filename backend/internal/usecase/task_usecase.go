package usecase

import (
	"todo-app/backend/internal/models"
	"todo-app/backend/internal/service"
)

type TaskUsecase interface {
	GetTasks(userID int64, filter string) ([]models.Task, error)
	AddTask(userID int64, title, desc, priority string) (*models.Task, error)
	UpdateTask(userID int64, task models.Task) (*models.Task, error)
	ToggleTask(userID, id int64) (*models.Task, error)
	DeleteTask(userID, id int64) error
	ClearCompleted(userID int64) (int64, error)
}

type taskUsecase struct {
	service service.TaskService
}

func NewTaskUsecase(s service.TaskService) TaskUsecase {
	return &taskUsecase{service: s}
}

func (u *taskUsecase) GetTasks(userID int64, filter string) ([]models.Task, error) {
	return u.service.GetTasks(userID, filter)
}

func (u *taskUsecase) AddTask(userID int64, title, desc, priority string) (*models.Task, error) {
	t := &models.Task{
		UserID:      userID,
		Title:       title,
		Description: desc,
		Priority:    priority,
	}
	return u.service.AddTask(t)
}

func (u *taskUsecase) UpdateTask(userID int64, task models.Task) (*models.Task, error) {
	task.UserID = userID
	return u.service.UpdateTask(&task)
}

func (u *taskUsecase) ToggleTask(userID, id int64) (*models.Task, error) {
	return u.service.ToggleTask(userID, id)
}

func (u *taskUsecase) DeleteTask(userID, id int64) error {
	return u.service.DeleteTask(userID, id)
}

func (u *taskUsecase) ClearCompleted(userID int64) (int64, error) {
	return u.service.ClearCompleted(userID)
}
