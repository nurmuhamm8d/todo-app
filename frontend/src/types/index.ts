export type Priority = 'high' | 'medium' | 'low';
export type Filter = 'all' | 'active' | 'completed';
export type SortOrder = 'date' | 'priority';

export interface TaskDTO {
  id: number;
  title: string;
  priority: Priority;
  completed: boolean;
  createdAt: string;
  completedAt: string | null;
  dueDate: string | Date | null;
  [key: string]: unknown; // Allow additional properties
}

export interface StatsDTO {
  total: number;
  active: number;
  completed: number;
}

export interface AppError {
  message: string;
  code?: string;
  details?: unknown;
}

export interface ThemeContextType {
  theme: 'light' | 'dark';
  toggleTheme: () => void;
}

export interface TaskFormData {
  title: string;
  priority: Priority;
  dueDate: Date | string | null;
}

export interface TaskFilters {
  status: Filter;
  priority?: Priority;
  search?: string;
}

export interface TaskSorting {
  field: 'createdAt' | 'priority' | 'dueDate' | 'title';
  order: 'asc' | 'desc';
}

export interface TaskStats {
  total: number;
  completed: number;
  active: number;
  byPriority: {
    high: number;
    medium: number;
    low: number;
  };
  overdue: number;
}
