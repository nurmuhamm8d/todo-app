package main

import "database/sql"

func ensureSchema(db *sql.DB) error {
	_, err := db.Exec(`
create table if not exists users (
  id bigserial primary key,
  email text unique not null,
  password_hash text not null,
  created_at timestamptz not null default now()
);
create table if not exists tasks (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,
  title text not null,
  priority text not null default 'medium',
  completed boolean not null default false,
  created_at timestamptz not null default now(),
  completed_at timestamptz,
  due_at timestamptz
);
create index if not exists idx_tasks_user on tasks(user_id);
`)
	return err
}
