package main

import "database/sql"

func ensureSchema(db *sql.DB) error {
	_, err := db.Exec(`
create table if not exists tasks (
  id bigserial primary key,
  user_id bigint not null,
  title text not null,
  priority text not null default 'medium',
  completed boolean not null default false,
  created_at timestamptz not null default now(),
  due_at timestamptz null,
  completed_at timestamptz null
);
alter table tasks add column if not exists description text default '';
alter table tasks add column if not exists tags text[] default '{}';
alter table tasks add column if not exists repeat_rule text;
alter table tasks add column if not exists category_id bigint;

create table if not exists categories (
  id bigserial primary key,
  user_id bigint not null,
  name text not null
);

create table if not exists subtasks (
  id bigserial primary key,
  user_id bigint not null,
  task_id bigint not null references tasks(id) on delete cascade,
  title text not null,
  completed boolean not null default false,
  created_at timestamptz not null default now()
);

create index if not exists idx_tasks_user on tasks(user_id);
create index if not exists idx_tasks_due on tasks(user_id, due_at);
create index if not exists idx_tasks_tags on tasks using gin (tags);
create index if not exists idx_subtasks_task on subtasks(user_id, task_id);
`)
	return err
}
