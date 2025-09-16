// Global type definitions for the application

// Extend the Wails namespace with our custom types
declare namespace wails {
  namespace go {
    namespace main {
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
