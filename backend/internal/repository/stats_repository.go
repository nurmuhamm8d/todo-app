package repository

import "database/sql"

type StatsSnapshot struct {
	Total          int
	Active         int
	Completed      int
	Overdue        int
	HighPriority   int
	MediumPriority int
	LowPriority    int
}

type StatsRepository interface {
	Snapshot() (StatsSnapshot, error)
}

type statsRepo struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) StatsRepository {
	return &statsRepo{db: db}
}

func (r *statsRepo) Snapshot() (StatsSnapshot, error) {
	var s StatsSnapshot
	err := r.db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&s.Total)
	if err != nil {
		return s, err
	}
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE completed=false`).Scan(&s.Active)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE completed=true`).Scan(&s.Completed)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE completed=false AND due_date < now()`).Scan(&s.Overdue)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE priority='high'`).Scan(&s.HighPriority)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE priority='medium'`).Scan(&s.MediumPriority)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE priority='low'`).Scan(&s.LowPriority)
	return s, nil
}
