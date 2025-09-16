package repository

import (
	"database/sql"

	"todo-app/backend/internal/models"
)

type CategoryRepository interface {
	Create(userID int64, name string) (models.Category, error)
	List(userID int64) ([]models.Category, error)
}

type categoryRepo struct{ db *sql.DB }

func NewCategoryRepository(db *sql.DB) CategoryRepository { return &categoryRepo{db: db} }

func (r *categoryRepo) Create(userID int64, name string) (models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(`INSERT INTO categories(user_id,name) VALUES($1,$2) RETURNING id,user_id,name`, userID, name).Scan(&c.ID, &c.UserID, &c.Name)
	return c, err
}

func (r *categoryRepo) List(userID int64) ([]models.Category, error) {
	rows, err := r.db.Query(`SELECT id,user_id,name FROM categories WHERE user_id=$1 ORDER BY name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}
