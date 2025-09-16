package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserPublic struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type Category struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"-"`
	Name   string `json:"name"`
}

type Task struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"-"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Priority     string     `json:"priority"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"createdAt"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	RepeatRule   *string    `json:"repeatRule,omitempty"`
	CategoryID   *int64     `json:"categoryId,omitempty"`
	CategoryName *string    `json:"categoryName,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
}

type CreateTaskInput struct {
	UserID      int64      `json:"-"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	RepeatRule  *string    `json:"repeatRule,omitempty"`
	CategoryID  *int64     `json:"categoryId,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

type UpdateTaskInput struct {
	UserID      int64      `json:"-"`
	ID          int64      `json:"id"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	RepeatRule  *string    `json:"repeatRule,omitempty"`
	CategoryID  *int64     `json:"categoryId,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

type TaskFilter struct {
	UserID     int64      `json:"-"`
	Status     string     `json:"status"`
	Priority   string     `json:"priority"`
	DateFilter string     `json:"dateFilter"`
	Search     string     `json:"search"`
	CategoryID *int64     `json:"categoryId,omitempty"`
	From       *time.Time `json:"from,omitempty"`
	To         *time.Time `json:"to,omitempty"`
	Sort       string     `json:"sort"`
}

type Stats struct {
	Total        int64 `json:"total"`
	Active       int64 `json:"active"`
	Completed    int64 `json:"completed"`
	Overdue      int64 `json:"overdue"`
	HighPriority int64 `json:"highPriority"`
}
