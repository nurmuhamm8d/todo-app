create table if not exists categories (
  id bigserial primary key,
  user_id bigint not null,
  name text not null,
  created_at timestamp not null default now()
);

create table if not exists tasks (
  id bigserial primary key,
  user_id bigint not null,
  title text not null,
  description text not null default '',
  priority text not null default 'medium',
  completed boolean not null default false,
  created_at timestamp not null default now(),
  due_at timestamp null,
  completed_at timestamp null,
  repeat_rule text null,
  category_id bigint null references categories(id) on delete set null,
  tags text[] not null default '{}'
);

create index if not exists idx_tasks_user on tasks(user_id);
create index if not exists idx_tasks_completed on tasks(user_id, completed);
create index if not exists idx_tasks_due on tasks(user_id, due_at);
create index if not exists idx_tasks_overdue on tasks(user_id, completed, due_at);
create index if not exists idx_tasks_tags_gin on tasks using gin (tags);

create table if not exists subtasks (
  id bigserial primary key,
  user_id bigint not null,
  task_id bigint not null references tasks(id) on delete cascade,
  title text not null,
  completed boolean not null default false,
  created_at timestamp not null default now()
);

create index if not exists idx_subtasks_task on subtasks(user_id, task_id);
