package repository

import (
	"database/sql"
	"errors"

	"todo-app/backend/internal/models"
)

type AuthRepository interface {
	CreateUser(email, passwordHash string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id int64) (models.User, error)
}

type authRepo struct{ db *sql.DB }

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(email, passwordHash string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, email, password_hash, created_at
	`, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

func (r *authRepo) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *authRepo) GetUserByID(id int64) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return u, nil
}
