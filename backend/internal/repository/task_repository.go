package repository

import (
	"database/sql"
	"todo-app/backend/internal/models"
)

type TaskRepository interface {
	GetAll(userID int64, filter string) ([]models.Task, error)
	GetByID(userID, id int64) (*models.Task, error)
	Create(task *models.Task) (*models.Task, error)
	Update(task *models.Task) (*models.Task, error)
	Delete(userID, id int64) error
	ClearCompleted(userID int64) (int64, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetAll(userID int64, filter string) ([]models.Task, error) {
	q := `
		select id, user_id, title, coalesce(description,''), priority, completed,
		       created_at, due_at, completed_at, repeat_rule, category_id, coalesce(tags, '{}')
		from tasks
		where user_id = $1
	`
	switch filter {
	case "active":
		q += " and completed=false"
	case "completed":
		q += " and completed=true"
	}
	q += " order by created_at desc"

	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description, &t.Priority, &t.Completed,
			&t.CreatedAt, &t.DueDate, &t.CompletedAt, &t.RepeatRule, &t.CategoryID, &t.Tags,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

func (r *taskRepository) GetByID(userID, id int64) (*models.Task, error) {
	q := `
		select id, user_id, title, coalesce(description,''), priority, completed,
		       created_at, due_at, completed_at, repeat_rule, category_id, coalesce(tags, '{}')
		from tasks where user_id=$1 and id=$2
	`
	var t models.Task
	err := r.db.QueryRow(q, userID, id).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.Priority, &t.Completed,
		&t.CreatedAt, &t.DueDate, &t.CompletedAt, &t.RepeatRule, &t.CategoryID, &t.Tags,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepository) Create(task *models.Task) (*models.Task, error) {
	q := `
		insert into tasks (user_id, title, description, priority, completed, created_at, due_at, repeat_rule, category_id, tags)
		values ($1,$2,$3,$4,false,now(),$5,$6,$7,$8)
		returning id, created_at
	`
	err := r.db.QueryRow(q, task.UserID, task.Title, task.Description, task.Priority,
		task.DueDate, task.RepeatRule, task.CategoryID, task.Tags).
		Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) Update(task *models.Task) (*models.Task, error) {
	q := `
		update tasks set title=$1, description=$2, priority=$3, due_at=$4, repeat_rule=$5, category_id=$6, tags=$7
		where id=$8 and user_id=$9
		returning created_at, completed, completed_at
	`
	err := r.db.QueryRow(q, task.Title, task.Description, task.Priority, task.DueDate,
		task.RepeatRule, task.CategoryID, task.Tags, task.ID, task.UserID).
		Scan(&task.CreatedAt, &task.Completed, &task.CompletedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) Delete(userID, id int64) error {
	_, err := r.db.Exec(`delete from tasks where id=$1 and user_id=$2`, id, userID)
	return err
}

func (r *taskRepository) ClearCompleted(userID int64) (int64, error) {
	res, err := r.db.Exec(`delete from tasks where user_id=$1 and completed=true`, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
