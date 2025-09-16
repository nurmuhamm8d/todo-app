// Type definitions for Wails Go bindings

// This file provides TypeScript type definitions for the Wails Go backend bindings

interface TaskDTO {
  id: number;
  title: string;
  priority: 'low' | 'medium' | 'high';
  completed: boolean;
  createdAt: string;
  completedAt: string | null;
  dueDate: string | null;
}

interface StatsDTO {
  total: number;
  active: number;
  completed: number;
}

declare module "wailsjs/go/main/App" {
  export interface App {
    GetTasks(filter: 'all' | 'active' | 'completed'): Promise<TaskDTO[]>;
    AddTask(title: string, priority: 'low' | 'medium' | 'high', dueISO: string): Promise<TaskDTO>;
    ToggleTask(id: number): Promise<TaskDTO>;
    DeleteTask(id: number): Promise<void>;
    ClearCompleted(): Promise<number>;
    UpdateTask(id: number, title: string, priority: 'low' | 'medium' | 'high', dueISO: string): Promise<TaskDTO>;
    GetStats(): Promise<StatsDTO>;
  }
}

// Global type augmentations
declare global {
  namespace wails {
    namespace go {
      namespace main {
        interface App {
          GetTasks(filter: 'all' | 'active' | 'completed'): Promise<TaskDTO[]>;
          AddTask(title: string, priority: 'low' | 'medium' | 'high', dueISO: string): Promise<TaskDTO>;
          ToggleTask(id: number): Promise<TaskDTO>;
          DeleteTask(id: number): Promise<void>;
          ClearCompleted(): Promise<number>;
          UpdateTask(id: number, title: string, priority: 'low' | 'medium' | 'high', dueISO: string): Promise<TaskDTO>;
          GetStats(): Promise<StatsDTO>;
        }
      }
    }
  }
}

export type { TaskDTO, StatsDTO };
