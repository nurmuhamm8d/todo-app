import React, { useEffect, useMemo, useState } from "react";
import "./index.css";
import { GetTasks, AddTask, ToggleTask, DeleteTask, ClearCompleted, GetStats } from "./wailsjs/go/main/App";
import type { main as NS } from "./wailsjs/go/models";

type TaskDTO = NS.TaskDTO;
type StatsDTO = NS.StatsDTO;
type Priority = "high" | "medium" | "low";
type Filter = "all" | "active" | "completed";
type SortOrder = "date" | "priority";

function asArray<T>(v: unknown): T[] {
  return Array.isArray(v) ? (v as T[]) : [];
}

export default function App(): JSX.Element {
  const [tasks, setTasks] = useState<TaskDTO[]>([]);
  const [stats, setStats] = useState<StatsDTO | null>(null);
  const [title, setTitle] = useState("");
  const [priority, setPriority] = useState<Priority>("medium");
  const [dueLocal, setDueLocal] = useState<string>("");
  const [filter, setFilter] = useState<Filter>("all");
  const [sort, setSort] = useState<SortOrder>("date");

  async function refresh() {
    try {
      const list = await GetTasks(filter);
      setTasks(asArray<TaskDTO>(list));
    } catch {
      setTasks([]);
    }
    try {
      const s = await GetStats();
      setStats(s as StatsDTO);
    } catch {
      setStats(null);
    }
  }

  useEffect(() => {
    refresh();
  }, [filter]);

  const sorted = useMemo(() => {
    const data = asArray<TaskDTO>(tasks).slice();
    if (sort === "date") {
      data.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
    } else {
      const order = { high: 0, medium: 1, low: 2 } as const;
      data.sort((a, b) => order[a.priority as keyof typeof order] - order[b.priority as keyof typeof order]);
    }
    return data;
  }, [tasks, sort]);

  function toRFC3339(localValue: string): string {
    if (!localValue) return "";
    const d = new Date(localValue);
    if (isNaN(d.getTime())) return "";
    return d.toISOString();
  }

  async function onAdd() {
    const tTitle = title.trim();
    if (!tTitle) return;
    try {
      const dueISO = toRFC3339(dueLocal);
      const t = await AddTask(tTitle, priority, dueISO);
      setTasks(prev => [t as TaskDTO, ...asArray<TaskDTO>(prev)]);
      setTitle("");
      setDueLocal("");
      const s = await GetStats();
      setStats(s as StatsDTO);
    } catch {}
  }

  async function onToggle(id: number) {
    try {
      const updated = await ToggleTask(id);
      setTasks(prev => asArray<TaskDTO>(prev).map(t => (t.id === (updated as TaskDTO).id ? (updated as TaskDTO) : t)));
      const s = await GetStats();
      setStats(s as StatsDTO);
    } catch {}
  }

  async function onDelete(id: number) {
    try {
      await DeleteTask(id);
      setTasks(prev => asArray<TaskDTO>(prev).filter(t => t.id !== id));
      const s = await GetStats();
      setStats(s as StatsDTO);
    } catch {}
  }

  async function onClearCompleted() {
    try {
      await ClearCompleted();
      const list = await GetTasks(filter);
      setTasks(asArray<TaskDTO>(list));
      const s = await GetStats();
      setStats(s as StatsDTO);
    } catch {
      setTasks([]);
    }
  }

  const total = stats?.total ?? 0;
  const active = stats?.active ?? 0;
  const completed = stats?.completed ?? 0;
  const overdue = stats?.overdue ?? 0;
  const progress = total > 0 ? Math.round((completed / total) * 100) : 0;

  return (
    <div className="container">
      <div className="shell">
        <div className="toolbar">
          <div className="brand">Todo App</div>
          <select className="select" value={filter} onChange={e => setFilter(e.target.value as Filter)}>
            <option value="all">All</option>
            <option value="active">Active</option>
            <option value="completed">Completed</option>
          </select>
          <select className="select" value={sort} onChange={e => setSort(e.target.value as SortOrder)}>
            <option value="date">By Date</option>
            <option value="priority">By Priority</option>
          </select>
          <button className="button secondary" onClick={onClearCompleted}>Clear Completed</button>
        </div>

        <div className="stats">
          <div className="stat"><div className="label">Total</div><div className="value">{total}</div></div>
          <div className="stat"><div className="label">Active</div><div className="value">{active}</div></div>
          <div className="stat"><div className="label">Completed</div><div className="value">{completed}</div></div>
          <div className="stat"><div className="label">Overdue</div><div className="value">{overdue}</div></div>
        </div>
        <div className="progress"><div style={{ width: `${progress}%` }} /></div>

        <div className="form">
          <input className="input" placeholder="Task title" value={title} onChange={e => setTitle(e.target.value)} />
          <select className="select" value={priority} onChange={e => setPriority(e.target.value as Priority)}>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
          <input className="input date" type="datetime-local" value={dueLocal} onChange={e => setDueLocal(e.target.value)} />
          <button className="button" onClick={onAdd}>Add</button>
        </div>

        <div className="filters">
          <div className={`chip ${filter==='all'?'active':''}`}>All</div>
          <div className={`chip ${filter==='active'?'active':''}`}>Active</div>
          <div className={`chip ${filter==='completed'?'active':''}`}>Completed</div>
        </div>

        <div className="list">
          {sorted.map(t => (
            <div key={t.id} className="card">
              <input className="checkbox" type="checkbox" checked={t.completed} onChange={() => onToggle(t.id)} />
              <div style={{flex:1}}>
                <div className="title">{t.title}</div>
                <div className="meta">{t.priority.toUpperCase()} • {new Date(t.createdAt).toLocaleString()}</div>
              </div>
              <div className={`badge ${t.priority}`}>{t.priority}</div>
              <div className="actions">
                <button className="icon-btn danger" onClick={() => onDelete(t.id)}>×</button>
              </div>
            </div>
          ))}
        </div>

        <div className="footer">
          <div>{progress}% complete</div>
          <div>{total} tasks</div>
        </div>
      </div>
    </div>
  );
}
