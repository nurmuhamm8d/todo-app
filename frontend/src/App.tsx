import React, { useEffect, useMemo, useState } from "react";
import "./index.css";
import ThemeToggle from "./components/ThemeToggle";
import Select from "./components/Select";
import { GetTasks, AddTask, ToggleTask, DeleteTask, ClearCompleted, GetStats } from "./wailsjs/go/main/App";
import type { main as NS } from "./wailsjs/go/models";

type TaskDTO = NS.TaskDTO;
type StatsDTO = NS.StatsDTO;
type Priority = "high" | "medium" | "low";
type Filter = "all" | "active" | "completed" | "overdue" | "today" | "week";
type SortOrder = "date" | "priority";

type Meta = { category: string; tags: string[] };
type MetaMap = Record<number, Meta>;

const META_KEY = "task-meta-v1";

function loadMeta(): MetaMap {
  try {
    const raw = localStorage.getItem(META_KEY);
    return raw ? JSON.parse(raw) : {};
  } catch {
    return {};
  }
}

function saveMeta(m: MetaMap) {
  localStorage.setItem(META_KEY, JSON.stringify(m));
}

export default function App(): JSX.Element {
  const [tasks, setTasks] = useState<TaskDTO[]>([]);
  const [stats, setStats] = useState<StatsDTO | null>(null);
  const [title, setTitle] = useState("");
  const [priority, setPriority] = useState<Priority>("medium");
  const [dueDate, setDueDate] = useState<string>("");
  const [filter, setFilter] = useState<Filter>("all");
  const [sort, setSort] = useState<SortOrder>("date");
  const [search, setSearch] = useState("");
  const [meta, setMeta] = useState<MetaMap>(() => loadMeta());
  const [newCategory, setNewCategory] = useState("");
  const [newTagsText, setNewTagsText] = useState("");
  const [selected, setSelected] = useState<Set<number>>(new Set());
  const [catFilter, setCatFilter] = useState<string>("all");

  async function refresh(currentFilter: Filter) {
    const list = await GetTasks(currentFilter);
    setTasks(Array.isArray(list) ? list : []);
    const s = await GetStats();
    setStats(s);
  }

  useEffect(() => {
    refresh(filter);
  }, [filter]);

  useEffect(() => {
    saveMeta(meta);
  }, [meta]);

  const categories = useMemo(() => {
    const set = new Set<string>();
    tasks.forEach(t => {
      const m = meta[t.id];
      if (m?.category) set.add(m.category);
    });
    return ["all", ...Array.from(set).sort((a, b) => a.localeCompare(b))];
  }, [tasks, meta]);

  const enriched = useMemo(() => {
    return tasks.map(t => {
      const m = meta[t.id] ?? { category: "General", tags: [] };
      return { ...t, __category: m.category || "General", __tags: m.tags || [] as string[] };
    });
  }, [tasks, meta]);

  const filtered = useMemo(() => {
    const s = search.trim().toLowerCase();
    const cf = catFilter;
    return enriched.filter(t => {
      const okTitle = s ? t.title.toLowerCase().includes(s) : true;
      const okCat = cf === "all" ? true : t.__category === cf;
      return okTitle && okCat;
    });
  }, [enriched, search, catFilter]);

  const sorted = useMemo(() => {
    const copy = [...filtered];
    if (sort === "date") {
      copy.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
    } else {
      const order = { high: 0, medium: 1, low: 2 } as const;
      copy.sort((a, b) => order[a.priority as keyof typeof order] - order[b.priority as keyof typeof order]);
    }
    return copy;
  }, [filtered, sort]);

  const toISO = (s: string) => (!s ? "" : new Date(s).toISOString());

  function parseTags(input: string) {
    return input
      .split(",")
      .map(t => t.trim())
      .filter(Boolean)
      .slice(0, 12);
  }

  async function onAdd() {
    const t = title.trim();
    if (!t) return;
    const created = await AddTask(t, priority, toISO(dueDate));
    setTasks(prev => [created, ...prev]);
    const cat = newCategory.trim() || "General";
    const tags = parseTags(newTagsText);
    setMeta(prev => ({ ...prev, [created.id]: { category: cat, tags } }));
    setTitle("");
    setDueDate("");
    setNewCategory("");
    setNewTagsText("");
    setStats(await GetStats());
  }

  async function onToggle(id: number) {
    const before = tasks.find(x => x.id === id);
    const updated = await ToggleTask(id);
    setTasks(prev => prev.map(t => (t.id === updated.id ? updated : t)));
    setStats(await GetStats());
    if (before && before.completed !== updated.completed) {
      setSelected(prev => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  }

  async function onDelete(id: number) {
    await DeleteTask(id);
    setMeta(prev => {
      const n = { ...prev };
      delete n[id];
      return n;
    });
    await refresh(filter);
    setSelected(prev => {
      const n = new Set(prev);
      n.delete(id);
      return n;
    });
  }

  async function onClearCompleted() {
    await ClearCompleted();
    const remainingMeta: MetaMap = {};
    tasks.forEach(t => {
      if (!t.completed && meta[t.id]) remainingMeta[t.id] = meta[t.id];
    });
    setMeta(remainingMeta);
    await refresh(filter);
    setSelected(new Set());
  }

  function toggleSelect(id: number) {
    setSelected(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  function selectAllVisible(v: boolean) {
    if (v) {
      setSelected(new Set(sorted.map(t => t.id)));
    } else {
      setSelected(new Set());
    }
  }

  async function completeSelected() {
    const ids = Array.from(selected);
    const need = tasks.filter(t => ids.includes(t.id) && !t.completed).map(t => t.id);
    for (const id of need) {
      await ToggleTask(id);
    }
    await refresh(filter);
    setSelected(new Set());
  }

  async function deleteSelected() {
    const ids = Array.from(selected);
    for (const id of ids) {
      await DeleteTask(id);
    }
    const m: MetaMap = { ...meta };
    ids.forEach(id => delete m[id]);
    setMeta(m);
    await refresh(filter);
    setSelected(new Set());
  }

  const total = stats?.total ?? 0;
  const active = stats?.active ?? 0;
  const completed = stats?.completed ?? 0;
  const overdue = stats?.overdue ?? 0;
  const progress = total > 0 ? Math.round((completed / total) * 100) : 0;

  const filterOptions = [
    { value: "all", label: "All Tasks" },
    { value: "active", label: "Active" },
    { value: "completed", label: "Completed" },
    { value: "overdue", label: "Overdue" },
    { value: "today", label: "Today" },
    { value: "week", label: "This Week" }
  ];
  const sortOptions = [
    { value: "date", label: "Sort by Date" },
    { value: "priority", label: "Sort by Priority" }
  ];
  const priorityOptions = [
    { value: "high", label: "High" },
    { value: "medium", label: "Medium" },
    { value: "low", label: "Low" }
  ];
  const catOptions = categories.map(c => ({ value: c, label: c === "all" ? "All Categories" : c }));

  const anySelected = selected.size > 0;
  const allVisibleSelected = sorted.length > 0 && selected.size === sorted.length;

  return (
    <div className="app">
      <div className="container">
        <header className="app-header">
          <h1 className="app-title">Task Manager</h1>
          <div className="header-controls">
            <Select value={filter} onChange={(v) => setFilter(v as Filter)} options={filterOptions} className="w-200" />
            <Select value={sort} onChange={(v) => setSort(v as SortOrder)} options={sortOptions} className="w-200" />
            <Select value={catFilter} onChange={(v) => setCatFilter(v)} options={catOptions} className="w-200" />
            <input className="search-input" placeholder="Search by title..." value={search} onChange={e => setSearch(e.target.value)} />
            <button className="btn btn-clear" onClick={onClearCompleted}>Clear Completed</button>
            <ThemeToggle />
          </div>
        </header>

        <div className="stats-grid">
          <div className="stat-card"><div className="stat-value">{total}</div><div className="stat-label">Total</div></div>
          <div className="stat-card"><div className="stat-value">{active}</div><div className="stat-label">Active</div></div>
          <div className="stat-card"><div className="stat-value">{completed}</div><div className="stat-label">Completed</div></div>
          <div className="stat-card"><div className="stat-value">{overdue}</div><div className="stat-label">Overdue</div></div>
          <div className="progress-container">
            <div className="progress-bar"><div className="progress-fill" style={{ width: `${progress}%` }} /></div>
            <div className="progress-text">{progress}% complete</div>
          </div>
        </div>

        <div className="task-form">
          <div className="form-grid form-grid-extended">
            <input
              type="text"
              className="task-input"
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="Task title"
              onKeyDown={e => e.key === "Enter" && onAdd()}
            />
            <Select value={priority} onChange={(v) => setPriority(v as Priority)} options={priorityOptions} />
            <input type="text" className="category-input" placeholder="Category" value={newCategory} onChange={e => setNewCategory(e.target.value)} />
            <input type="text" className="tags-input" placeholder="tags, comma,separated" value={newTagsText} onChange={e => setNewTagsText(e.target.value)} />
            <input type="datetime-local" className="date-input" value={dueDate} onChange={e => setDueDate(e.target.value)} />
            <button className="btn btn-primary" onClick={onAdd} disabled={!title.trim()}>Add Task</button>
          </div>
        </div>

        <div className={`bulkbar ${anySelected ? "show" : ""}`}>
          <div className="bulk-left">
            <input
              type="checkbox"
              className="select-all"
              checked={allVisibleSelected}
              onChange={e => selectAllVisible(e.target.checked)}
            />
            <span className="bulk-text">{selected.size} selected</span>
          </div>
          <div className="bulk-actions">
            <button className="btn bulk-complete" onClick={completeSelected} disabled={!anySelected}>Complete Selected</button>
            <button className="btn bulk-delete" onClick={deleteSelected} disabled={!anySelected}>Delete Selected</button>
          </div>
        </div>

        <div className="task-list">
          {sorted.length > 0 ? (
            <ul className="task-items">
              {sorted.map(task => (
                <li key={task.id} className={`task-item ${task.completed ? "completed" : ""}`}>
                  <div className="task-content">
                    <input
                      type="checkbox"
                      className="select-checkbox"
                      checked={selected.has(task.id)}
                      onChange={() => toggleSelect(task.id)}
                    />
                    <input
                      type="checkbox"
                      className="task-checkbox"
                      checked={task.completed}
                      onChange={() => onToggle(task.id)}
                    />
                    <span className="task-title">{task.title}</span>
                    <span className={`priority-badge ${task.priority}`}>{task.priority}</span>
                    <span className="category-badge">{(meta[task.id]?.category) || "General"}</span>
                    {((meta[task.id]?.tags) || []).length > 0 && (
                      <span className="tags">
                        {meta[task.id]!.tags.map((tg, i) => (
                          <span key={tg + i} className="tag-chip">{tg}</span>
                        ))}
                      </span>
                    )}
                    {task.dueDate && <span className="due-date">{new Date(task.dueDate).toLocaleDateString()}</span>}
                  </div>
                  <button className="delete-btn" onClick={() => onDelete(task.id)} aria-label="Delete task">×</button>
                </li>
              ))}
            </ul>
          ) : (
            <div className="empty-state"><p>No tasks found. Add a new task to get started!</p></div>
          )}
        </div>
      </div>
    </div>
  );
}
