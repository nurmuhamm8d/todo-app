package repository

import (
	"database/sql"
	"time"

	"todo-app/backend/internal/models"

	"github.com/lib/pq"
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
		var (
			id          int64
			uID         int64
			title       string
			desc        string
			priority    string
			completed   bool
			createdAt   time.Time
			dueAt       sql.NullTime
			completedAt sql.NullTime
			repeatRule  sql.NullString
			categoryID  sql.NullInt64
			tags        pq.StringArray
		)
		if err := rows.Scan(
			&id, &uID, &title, &desc, &priority, &completed,
			&createdAt, &dueAt, &completedAt, &repeatRule, &categoryID, &tags,
		); err != nil {
			return nil, err
		}
		var t models.Task
		t.ID = id
		t.UserID = uID
		t.Title = title
		t.Description = desc
		t.Priority = priority
		t.Completed = completed
		t.CreatedAt = createdAt
		if dueAt.Valid {
			v := dueAt.Time
			t.DueDate = &v
		}
		if completedAt.Valid {
			v := completedAt.Time
			t.CompletedAt = &v
		}
		if repeatRule.Valid {
			v := repeatRule.String
			t.RepeatRule = &v
		}
		if categoryID.Valid {
			v := categoryID.Int64
			t.CategoryID = &v
		}
		t.Tags = []string(tags)
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
	var (
		tid         int64
		uID         int64
		title       string
		desc        string
		priority    string
		completed   bool
		createdAt   time.Time
		dueAt       sql.NullTime
		completedAt sql.NullTime
		repeatRule  sql.NullString
		categoryID  sql.NullInt64
		tags        pq.StringArray
	)
	if err := r.db.QueryRow(q, userID, id).Scan(
		&tid, &uID, &title, &desc, &priority, &completed,
		&createdAt, &dueAt, &completedAt, &repeatRule, &categoryID, &tags,
	); err != nil {
		return nil, err
	}
	var t models.Task
	t.ID = tid
	t.UserID = uID
	t.Title = title
	t.Description = desc
	t.Priority = priority
	t.Completed = completed
	t.CreatedAt = createdAt
	if dueAt.Valid {
		v := dueAt.Time
		t.DueDate = &v
	}
	if completedAt.Valid {
		v := completedAt.Time
		t.CompletedAt = &v
	}
	if repeatRule.Valid {
		v := repeatRule.String
		t.RepeatRule = &v
	}
	if categoryID.Valid {
		v := categoryID.Int64
		t.CategoryID = &v
	}
	t.Tags = []string(tags)
	return &t, nil
}

func (r *taskRepository) Create(task *models.Task) (*models.Task, error) {
	q := `
		insert into tasks (user_id, title, description, priority, completed, created_at, due_at, repeat_rule, category_id, tags)
		values ($1,$2,$3,$4,false,now(),$5,$6,$7,$8)
		returning id, created_at
	`
	if err := r.db.QueryRow(q,
		task.UserID,
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
		task.RepeatRule,
		task.CategoryID,
		pq.Array(task.Tags),
	).Scan(&task.ID, &task.CreatedAt); err != nil {
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
	var completedAt sql.NullTime
	if err := r.db.QueryRow(q,
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
		task.RepeatRule,
		task.CategoryID,
		pq.Array(task.Tags),
		task.ID,
		task.UserID,
	).Scan(&task.CreatedAt, &task.Completed, &completedAt); err != nil {
		return nil, err
	}
	if completedAt.Valid {
		v := completedAt.Time
		task.CompletedAt = &v
	} else {
		task.CompletedAt = nil
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
