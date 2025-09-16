package repository

import (
	"database/sql"
	"todo-app/backend/internal/models"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	Create(user *models.User) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`select id,email,password_hash,created_at from users where email=$1`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`select id,email,password_hash,created_at from users where id=$1`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	err := r.db.QueryRow(
		`insert into users (email,password_hash,created_at) values ($1,$2,now()) returning id,created_at`,
		user.Email, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
