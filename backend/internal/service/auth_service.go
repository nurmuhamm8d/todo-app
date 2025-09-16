package service

import (
	"errors"

	"todo-app/backend/internal/models"
	"todo-app/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password string) (models.UserPublic, error)
	Login(email, password string) (models.UserPublic, error)
	GetProfile(userID int64) (models.UserPublic, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) Register(email, password string) (models.UserPublic, error) {
	existing, _ := s.repo.GetUserByEmail(email)
	if existing.ID != 0 {
		return models.UserPublic{}, errors.New("user already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserPublic{}, err
	}
	u, err := s.repo.CreateUser(email, string(hash))
	if err != nil {
		return models.UserPublic{}, err
	}
	return models.UserPublic{ID: u.ID, Email: u.Email}, nil
}

func (s *authService) Login(email, password string) (models.UserPublic, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return models.UserPublic{}, err
	}
	if u.ID == 0 {
		return models.UserPublic{}, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return models.UserPublic{}, errors.New("invalid email or password")
	}
	return models.UserPublic{ID: u.ID, Email: u.Email}, nil
}

func (s *authService) GetProfile(userID int64) (models.UserPublic, error) {
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		return models.UserPublic{}, err
	}
	if u.ID == 0 {
		return models.UserPublic{}, errors.New("user not found")
	}
	return models.UserPublic{ID: u.ID, Email: u.Email}, nil
}
