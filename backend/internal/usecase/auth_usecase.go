package usecase

import (
	"todo-app/backend/internal/models"
	"todo-app/backend/internal/service"
)

type AuthUsecase interface {
	Register(email, password string) (models.UserPublic, error)
	Login(email, password string) (models.UserPublic, error)
	GetProfile(userID int64) (models.UserPublic, error)
}

type authUsecase struct {
	svc service.AuthService
}

func NewAuthUsecase(s service.AuthService) AuthUsecase {
	return &authUsecase{svc: s}
}

func (u *authUsecase) Register(email, password string) (models.UserPublic, error) {
	return u.svc.Register(email, password)
}

func (u *authUsecase) Login(email, password string) (models.UserPublic, error) {
	return u.svc.Login(email, password)
}

func (u *authUsecase) GetProfile(userID int64) (models.UserPublic, error) {
	return u.svc.GetProfile(userID)
}
