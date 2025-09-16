package repository

import "database/sql"

type TagRepository interface {
	All() ([]string, error)
}

type tagRepo struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepo{db: db}
}

func (r *tagRepo) All() ([]string, error) {
	rows, err := r.db.Query(`SELECT name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
