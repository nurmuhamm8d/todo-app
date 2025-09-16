package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	db  *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{db: db}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(a.ctx, "App started")
}

type TaskDTO struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Priority    string  `json:"priority"`
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"createdAt"`
	CompletedAt *string `json:"completedAt,omitempty"`
	DueDate     *string `json:"dueDate,omitempty"`
}

type StatsDTO struct {
	Total     int64 `json:"total"`
	Active    int64 `json:"active"`
	Completed int64 `json:"completed"`
	Overdue   int64 `json:"overdue"`
}

func sPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	v := t.UTC().Format(time.RFC3339)
	return &v
}

func parseRFC3339OrNil(s string) (*time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
		loc := time.Now().Location()
		tt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
		return &tt, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		loc := time.Now().Location()
		tt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		return &tt, nil
	}
	return nil, errors.New("invalid date format")
}

const userID int64 = 1

func (a *App) GetTasks(filter string) ([]TaskDTO, error) {
	q := `
select id, title, priority, completed, created_at, completed_at, due_at
from tasks
where user_id = $1
`
	switch filter {
	case "active":
		q += " and completed = false"
	case "completed":
		q += " and completed = true"
	case "overdue":
		q += " and completed = false and due_at is not null and due_at < now()"
	case "today":
		q += " and due_at >= date_trunc('day', now()) and due_at < date_trunc('day', now()) + interval '1 day'"
	case "week":
		q += " and due_at >= date_trunc('week', now()) and due_at < date_trunc('week', now()) + interval '1 week'"
	}
	q += " order by created_at desc, id desc"

	rows, err := a.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []TaskDTO
	for rows.Next() {
		var id int64
		var title, priority string
		var completed bool
		var createdAt time.Time
		var completedAt sql.NullTime
		var dueAt sql.NullTime
		if err := rows.Scan(&id, &title, &priority, &completed, &createdAt, &completedAt, &dueAt); err != nil {
			return nil, err
		}
		res = append(res, TaskDTO{
			ID:        id,
			Title:     title,
			Priority:  priority,
			Completed: completed,
			CreatedAt: createdAt.UTC().Format(time.RFC3339),
			CompletedAt: func() *string {
				if completedAt.Valid {
					return sPtr(&completedAt.Time)
				}
				return nil
			}(),
			DueDate: func() *string {
				if dueAt.Valid {
					return sPtr(&dueAt.Time)
				}
				return nil
			}(),
		})
	}
	return res, nil
}

func (a *App) AddTask(title, priority, dueISO string) (TaskDTO, error) {
	if strings.TrimSpace(title) == "" {
		return TaskDTO{}, errors.New("title is required")
	}
	switch priority {
	case "high", "medium", "low":
	default:
		priority = "medium"
	}
	due, err := parseRFC3339OrNil(dueISO)
	if err != nil {
		return TaskDTO{}, err
	}
	var id int64
	var retTitle, retPriority string
	var completed bool
	var createdAt time.Time
	var completedAt sql.NullTime
	var dueAt sql.NullTime
	err = a.db.QueryRow(`
insert into tasks (user_id, title, priority, completed, created_at, due_at)
values ($1, $2, $3, false, now(), $4)
returning id, title, priority, completed, created_at, completed_at, due_at
`, userID, strings.TrimSpace(title), priority, due).
		Scan(&id, &retTitle, &retPriority, &completed, &createdAt, &completedAt, &dueAt)
	if err != nil {
		return TaskDTO{}, err
	}
	return TaskDTO{
		ID:        id,
		Title:     retTitle,
		Priority:  retPriority,
		Completed: completed,
		CreatedAt: createdAt.UTC().Format(time.RFC3339),
		CompletedAt: func() *string {
			if completedAt.Valid {
				return sPtr(&completedAt.Time)
			}
			return nil
		}(),
		DueDate: func() *string {
			if dueAt.Valid {
				return sPtr(&dueAt.Time)
			}
			return nil
		}(),
	}, nil
}

func (a *App) ToggleTask(id int64) (TaskDTO, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback()

	var completed bool
	if err := tx.QueryRow(`select completed from tasks where id=$1 and user_id=$2`, id, userID).Scan(&completed); err != nil {
		return TaskDTO{}, err
	}
	if completed {
		if _, err := tx.Exec(`update tasks set completed=false, completed_at=null where id=$1 and user_id=$2`, id, userID); err != nil {
			return TaskDTO{}, err
		}
	} else {
		if _, err := tx.Exec(`update tasks set completed=true, completed_at=now() where id=$1 and user_id=$2`, id, userID); err != nil {
			return TaskDTO{}, err
		}
	}

	var retID int64
	var title, priority string
	var retCompleted bool
	var createdAt time.Time
	var completedAt sql.NullTime
	var dueAt sql.NullTime

	if err := tx.QueryRow(`
select id, title, priority, completed, created_at, completed_at, due_at
from tasks where id=$1 and user_id=$2
`, id, userID).Scan(&retID, &title, &priority, &retCompleted, &createdAt, &completedAt, &dueAt); err != nil {
		return TaskDTO{}, err
	}

	if err := tx.Commit(); err != nil {
		return TaskDTO{}, err
	}

	return TaskDTO{
		ID:        retID,
		Title:     title,
		Priority:  priority,
		Completed: retCompleted,
		CreatedAt: createdAt.UTC().Format(time.RFC3339),
		CompletedAt: func() *string {
			if completedAt.Valid {
				return sPtr(&completedAt.Time)
			}
			return nil
		}(),
		DueDate: func() *string {
			if dueAt.Valid {
				return sPtr(&dueAt.Time)
			}
			return nil
		}(),
	}, nil
}

func (a *App) DeleteTask(id int64) error {
	res, err := runtime.MessageDialog(a.ctx, dialog.MessageDialog{
		Type:          dialog.Question,
		Title:         "Confirm deletion",
		Message:       "Delete this task?",
		Buttons:       []string{"Delete", "Cancel"},
		DefaultButton: "Delete",
		CancelButton:  "Cancel",
	})
	if err != nil {
		return err
	}
	if res != "Delete" {
		return nil
	}
	_, err = a.db.Exec(`delete from tasks where id=$1 and user_id=$2`, id, userID)
	return err
}

func (a *App) ClearCompleted() (int64, error) {
	res, err := a.db.Exec(`delete from tasks where user_id=$1 and completed=true`, userID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (a *App) UpdateTask(id int64, title, priority, dueISO string) (TaskDTO, error) {
	if strings.TrimSpace(title) == "" {
		return TaskDTO{}, errors.New("title is required")
	}
	switch priority {
	case "high", "medium", "low":
	default:
		priority = "medium"
	}
	due, err := parseRFC3339OrNil(dueISO)
	if err != nil {
		return TaskDTO{}, err
	}
	_, err = a.db.Exec(`update tasks set title=$1, priority=$2, due_at=$3 where id=$4 and user_id=$5`,
		strings.TrimSpace(title), priority, due, id, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	var retID int64
	var retTitle, retPriority string
	var completed bool
	var createdAt time.Time
	var completedAt sql.NullTime
	var dueAt sql.NullTime
	err = a.db.QueryRow(`
select id, title, priority, completed, created_at, completed_at, due_at
from tasks where id=$1 and user_id=$2
`, id, userID).Scan(&retID, &retTitle, &retPriority, &completed, &createdAt, &completedAt, &dueAt)
	if err != nil {
		return TaskDTO{}, err
	}
	return TaskDTO{
		ID:        retID,
		Title:     retTitle,
		Priority:  retPriority,
		Completed: completed,
		CreatedAt: createdAt.UTC().Format(time.RFC3339),
		CompletedAt: func() *string {
			if completedAt.Valid {
				return sPtr(&completedAt.Time)
			}
			return nil
		}(),
		DueDate: func() *string {
			if dueAt.Valid {
				return sPtr(&dueAt.Time)
			}
			return nil
		}(),
	}, nil
}

func (a *App) GetStats() (StatsDTO, error) {
	var s StatsDTO
	if err := a.db.QueryRow(`select count(*) from tasks where user_id=$1`, userID).Scan(&s.Total); err != nil {
		return s, err
	}
	if err := a.db.QueryRow(`select count(*) from tasks where user_id=$1 and completed=false`, userID).Scan(&s.Active); err != nil {
		return s, err
	}
	if err := a.db.QueryRow(`select count(*) from tasks where user_id=$1 and completed=true`, userID).Scan(&s.Completed); err != nil {
		return s, err
	}
	if err := a.db.QueryRow(`select count(*) from tasks where user_id=$1 and completed=false and due_at is not null and due_at < now()`, userID).Scan(&s.Overdue); err != nil {
		return s, err
	}
	return s, nil
}
