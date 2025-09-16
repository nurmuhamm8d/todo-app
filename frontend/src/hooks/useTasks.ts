import { useState, useEffect, useCallback } from 'react';
import { TaskDTO, Filter, TaskFormData, AppError, Priority } from '../types';
import * as App from '../wailsjs/go/main/App';

export const useTasks = () => {
  const [tasks, setTasks] = useState<TaskDTO[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<AppError | null>(null);
  const [stats, setStats] = useState<{ total: number; active: number; completed: number }>({ 
    total: 0, 
    active: 0, 
    completed: 0 
  });

  const fetchTasks = useCallback(async (filter: Filter = 'all') => {
    setLoading(true);
    setError(null);
    try {
      const [tasks, stats] = await Promise.all([
        App.GetTasks(filter) as Promise<TaskDTO[]>,
        App.GetStats() as Promise<{ total: number; active: number; completed: number }>
      ]);
      setTasks(tasks);
      setStats(stats);
    } catch (err) {
      setError({
        message: 'Failed to fetch tasks',
        code: 'FETCH_ERROR',
        details: err
      });
    } finally {
      setLoading(false);
    }
  }, []);

  const addTask = useCallback(async (data: TaskFormData) => {
    setLoading(true);
    setError(null);
    try {
      // Ensure dueDate is properly formatted as ISO string or empty string
      const dueDate = data.dueDate 
        ? data.dueDate instanceof Date 
          ? data.dueDate.toISOString() 
          : new Date(data.dueDate).toISOString()
        : '';
      
      const newTask = await App.AddTask(
        data.title,
        data.priority,
        dueDate
      ) as TaskDTO;
      await fetchTasks();
    } catch (err) {
      setError({
        message: 'Failed to add task',
        code: 'ADD_ERROR',
        details: err
      });
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchTasks]);

  const toggleTask = useCallback(async (id: number) => {
    setLoading(true);
    setError(null);
    try {
      const updatedTask = await App.ToggleTask(id) as TaskDTO;
      await fetchTasks();
    } catch (err) {
      setError({
        message: 'Failed to toggle task status',
        code: 'TOGGLE_ERROR',
        details: err
      });
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchTasks]);

  const removeTask = useCallback(async (id: number) => {
    setLoading(true);
    setError(null);
    try {
      await App.DeleteTask(id);
      await fetchTasks();
    } catch (err) {
      setError({
        message: 'Failed to delete task',
        code: 'DELETE_ERROR',
        details: err
      });
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchTasks]);

  const clearCompletedTasks = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const count = await App.ClearCompleted() as number;
      await fetchTasks();
    } catch (err) {
      setError({
        message: 'Failed to clear completed tasks',
        code: 'CLEAR_ERROR',
        details: err
      });
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchTasks]);

  const editTask = useCallback(async (id: number, updates: Partial<TaskFormData>) => {
    setLoading(true);
    setError(null);
    try {
      // Ensure dueDate is properly formatted as ISO string or empty string
      const dueDate = updates.dueDate 
        ? updates.dueDate instanceof Date 
          ? updates.dueDate.toISOString() 
          : new Date(updates.dueDate).toISOString()
        : '';
      
      const updatedTask = await App.UpdateTask(
        id,
        updates.title || '',
        updates.priority || 'medium',
        dueDate
      ) as TaskDTO;
      await fetchTasks();
    } catch (err) {
      setError({
        message: 'Failed to update task',
        code: 'UPDATE_ERROR',
        details: err
      });
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchTasks]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  return {
    tasks,
    stats,
    loading,
    error,
    addTask,
    toggleTask,
    removeTask,
    clearCompletedTasks,
    editTask,
    refetch: fetchTasks
  };
};
