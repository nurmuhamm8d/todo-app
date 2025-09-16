package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
}

type TaskDTO struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Priority    string   `json:"priority"`
	Completed   bool     `json:"completed"`
	CreatedAt   string   `json:"createdAt"`
	CompletedAt *string  `json:"completedAt,omitempty"`
	DueDate     *string  `json:"dueDate,omitempty"`
	CategoryID  *int64   `json:"categoryId,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type SubtaskDTO struct {
	ID        int64  `json:"id"`
	TaskID    int64  `json:"taskId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"createdAt"`
}

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
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
		u := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
		return &u, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		loc := time.Now().Location()
		u := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		return &u, nil
	}
	return nil, errors.New("invalid date format")
}

const userID int64 = 1

func scanTask(rows *sql.Rows) (TaskDTO, error) {
	var id int64
	var title, priority string
	var completed bool
	var createdAt time.Time
	var completedAt sql.NullTime
	var dueAt sql.NullTime
	var categoryID sql.NullInt64
	var tags []sql.NullString
	err := rows.Scan(&id, &title, &priority, &completed, &createdAt, &completedAt, &dueAt, &categoryID, pq.Array(&tags))
	if err != nil {
		return TaskDTO{}, err
	}
	var t []string
	for _, v := range tags {
		if v.Valid {
			t = append(t, v.String)
		}
	}
	var cid *int64
	if categoryID.Valid {
		v := categoryID.Int64
		cid = &v
	}
	return TaskDTO{
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
		CategoryID: cid,
		Tags:       t,
	}, nil
}

func (a *App) GetTasks(filter string) ([]TaskDTO, error) {
	q := `
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
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
		q += " and due_at::date = current_date"
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
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

func (a *App) SearchTasks(query string) ([]TaskDTO, error) {
	q := `
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
from tasks
where user_id=$1 and (title ilike $2 or coalesce(description,'') ilike $2)
order by created_at desc, id desc
`
	rows, err := a.db.Query(q, userID, "%"+strings.TrimSpace(query)+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []TaskDTO
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
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
	var categoryID sql.NullInt64
	var tags []sql.NullString
	err = a.db.QueryRow(`
insert into tasks (user_id, title, priority, completed, created_at, due_at, tags)
values ($1,$2,$3,false,now(),$4,$5)
returning id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
`, userID, strings.TrimSpace(title), priority, due, pq.Array([]string{})).Scan(&id, &retTitle, &retPriority, &completed, &createdAt, &completedAt, &dueAt, &categoryID, pq.Array(&tags))
	if err != nil {
		return TaskDTO{}, err
	}
	var t []string
	for _, v := range tags {
		if v.Valid {
			t = append(t, v.String)
		}
	}
	var cid *int64
	if categoryID.Valid {
		v := categoryID.Int64
		cid = &v
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
		CategoryID: cid,
		Tags:       t,
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
		var rr sql.NullString
		var due sql.NullTime
		var ttl, prio string
		if err := tx.QueryRow(`select title, priority, repeat_rule, due_at from tasks where id=$1 and user_id=$2`, id, userID).Scan(&ttl, &prio, &rr, &due); err == nil {
			if rr.Valid && due.Valid {
				next := due.Time
				switch strings.ToLower(rr.String) {
				case "daily":
					next = next.Add(24 * time.Hour)
				case "weekly":
					next = next.AddDate(0, 0, 7)
				case "monthly":
					next = next.AddDate(0, 1, 0)
				default:
				}
				if _, err := tx.Exec(`insert into tasks (user_id, title, priority, completed, created_at, due_at, repeat_rule, tags, category_id) select user_id, $1, $2, false, now(), $3, repeat_rule, tags, category_id from tasks where id=$4`, ttl, prio, next, id); err != nil {
					return TaskDTO{}, err
				}
			}
		}
	}
	row := tx.QueryRow(`
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
from tasks where id=$1 and user_id=$2
`, id, userID)
	var r TaskDTO
	{
		var rid int64
		var title, priority string
		var rcompleted bool
		var createdAt time.Time
		var completedAt sql.NullTime
		var dueAt sql.NullTime
		var categoryID sql.NullInt64
		var tags []sql.NullString
		if err := row.Scan(&rid, &title, &priority, &rcompleted, &createdAt, &completedAt, &dueAt, &categoryID, pq.Array(&tags)); err != nil {
			return TaskDTO{}, err
		}
		var t []string
		for _, v := range tags {
			if v.Valid {
				t = append(t, v.String)
			}
		}
		var cid *int64
		if categoryID.Valid {
			v := categoryID.Int64
			cid = &v
		}
		r = TaskDTO{
			ID:        rid,
			Title:     title,
			Priority:  priority,
			Completed: rcompleted,
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
			CategoryID: cid,
			Tags:       t,
		}
	}
	if err := tx.Commit(); err != nil {
		return TaskDTO{}, err
	}
	return r, nil
}

func (a *App) DeleteTask(id int64) error {
	sel, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    "question",
		Title:   "Delete task",
		Message: "Delete this task permanently?",
	})
	if err != nil {
		return err
	}
	if strings.ToLower(sel) != "yes" && strings.ToLower(sel) != "ok" {
		return nil
	}
	_, err = a.db.Exec(`delete from tasks where id=$1 and user_id=$2`, id, userID)
	return err
}

func (a *App) ClearCompleted() (int64, error) {
	sel, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    "question",
		Title:   "Clear completed",
		Message: "Delete all completed tasks?",
	})
	if err != nil {
		return 0, err
	}
	if strings.ToLower(sel) != "yes" && strings.ToLower(sel) != "ok" {
		return 0, nil
	}
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
	_, err = a.db.Exec(`update tasks set title=$1, priority=$2, due_at=$3 where id=$4 and user_id=$5`, strings.TrimSpace(title), priority, due, id, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	rows, err := a.db.Query(`
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
from tasks where id=$1 and user_id=$2
`, id, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	defer rows.Close()
	if rows.Next() {
		return scanTask(rows)
	}
	return TaskDTO{}, errors.New("not found")
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

func (a *App) SetTaskTags(id int64, tags []string) (TaskDTO, error) {
	_, err := a.db.Exec(`update tasks set tags=$1 where id=$2 and user_id=$3`, pq.Array(tags), id, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	rows, err := a.db.Query(`
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
from tasks where id=$1 and user_id=$2
`, id, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	defer rows.Close()
	if rows.Next() {
		return scanTask(rows)
	}
	return TaskDTO{}, errors.New("not found")
}

func (a *App) GetCategories() ([]CategoryDTO, error) {
	rows, err := a.db.Query(`select id, name from categories where user_id=$1 order by name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []CategoryDTO
	for rows.Next() {
		var c CategoryDTO
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (a *App) AddCategory(name string) (CategoryDTO, error) {
	if strings.TrimSpace(name) == "" {
		return CategoryDTO{}, errors.New("name is required")
	}
	var id int64
	if err := a.db.QueryRow(`insert into categories (user_id, name, created_at) values ($1,$2,now()) returning id`, userID, strings.TrimSpace(name)).Scan(&id); err != nil {
		return CategoryDTO{}, err
	}
	return CategoryDTO{ID: id, Name: strings.TrimSpace(name)}, nil
}

func (a *App) DeleteCategory(id int64) error {
	_, err := a.db.Exec(`delete from categories where id=$1 and user_id=$2`, id, userID)
	return err
}

func (a *App) AssignCategory(taskID, categoryID int64) (TaskDTO, error) {
	_, err := a.db.Exec(`update tasks set category_id=$1 where id=$2 and user_id=$3`, categoryID, taskID, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	rows, err := a.db.Query(`
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
from tasks where id=$1 and user_id=$2
`, taskID, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	defer rows.Close()
	if rows.Next() {
		return scanTask(rows)
	}
	return TaskDTO{}, errors.New("not found")
}

func (a *App) ClearCategory(taskID int64) (TaskDTO, error) {
	_, err := a.db.Exec(`update tasks set category_id=null where id=$1 and user_id=$2`, taskID, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	rows, err := a.db.Query(`
select id, title, priority, completed, created_at, completed_at, due_at, category_id, tags
from tasks where id=$1 and user_id=$2
`, taskID, userID)
	if err != nil {
		return TaskDTO{}, err
	}
	defer rows.Close()
	if rows.Next() {
		return scanTask(rows)
	}
	return TaskDTO{}, errors.New("not found")
}

func (a *App) GetSubtasks(taskID int64) ([]SubtaskDTO, error) {
	rows, err := a.db.Query(`select id, task_id, title, completed, created_at from subtasks where user_id=$1 and task_id=$2 order by id`, userID, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []SubtaskDTO
	for rows.Next() {
		var s SubtaskDTO
		var created time.Time
		if err := rows.Scan(&s.ID, &s.TaskID, &s.Title, &s.Completed, &created); err != nil {
			return nil, err
		}
		s.CreatedAt = created.UTC().Format(time.RFC3339)
		res = append(res, s)
	}
	return res, nil
}

func (a *App) AddSubtask(taskID int64, title string) (SubtaskDTO, error) {
	if strings.TrimSpace(title) == "" {
		return SubtaskDTO{}, errors.New("title is required")
	}
	var id int64
	var created time.Time
	if err := a.db.QueryRow(`insert into subtasks (user_id, task_id, title, completed, created_at) values ($1,$2,$3,false,now()) returning id, created_at`, userID, taskID, strings.TrimSpace(title)).Scan(&id, &created); err != nil {
		return SubtaskDTO{}, err
	}
	return SubtaskDTO{ID: id, TaskID: taskID, Title: strings.TrimSpace(title), Completed: false, CreatedAt: created.UTC().Format(time.RFC3339)}, nil
}

func (a *App) ToggleSubtask(id int64) (SubtaskDTO, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return SubtaskDTO{}, err
	}
	defer tx.Rollback()
	var completed bool
	var taskID int64
	if err := tx.QueryRow(`select completed, task_id from subtasks where id=$1 and user_id=$2`, id, userID).Scan(&completed, &taskID); err != nil {
		return SubtaskDTO{}, err
	}
	if _, err := tx.Exec(`update subtasks set completed=$1 where id=$2 and user_id=$3`, !completed, id, userID); err != nil {
		return SubtaskDTO{}, err
	}
	var s SubtaskDTO
	var created time.Time
	if err := tx.QueryRow(`select id, task_id, title, completed, created_at from subtasks where id=$1`, id).Scan(&s.ID, &s.TaskID, &s.Title, &s.Completed, &created); err != nil {
		return SubtaskDTO{}, err
	}
	s.CreatedAt = created.UTC().Format(time.RFC3339)
	if err := tx.Commit(); err != nil {
		return SubtaskDTO{}, err
	}
	return s, nil
}

func (a *App) DeleteSubtask(id int64) error {
	_, err := a.db.Exec(`delete from subtasks where id=$1 and user_id=$2`, id, userID)
	return err
}

func (a *App) BulkComplete(ids []int64) (int64, error) {
	res, err := a.db.Exec(`update tasks set completed=true, completed_at=now() where user_id=$1 and completed=false and id = any($2)`, userID, pq.Array(ids))
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (a *App) BulkDelete(ids []int64) (int64, error) {
	sel, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    "question",
		Title:   "Delete selected",
		Message: "Delete selected tasks?",
	})
	if err != nil {
		return 0, err
	}
	if strings.ToLower(sel) != "yes" && strings.ToLower(sel) != "ok" {
		return 0, nil
	}
	res, err := a.db.Exec(`delete from tasks where user_id=$1 and id = any($2)`, userID, pq.Array(ids))
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
