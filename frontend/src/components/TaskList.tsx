import React, { useState } from 'react';
import { TaskItem } from './TaskItem';
import { TaskForm } from './TaskForm';
import { TaskDTO } from '../types';
import { Priority, Filter } from '../types';

// Helper function to convert date string to Date object
const parseDate = (dateString: string | Date | null | undefined): Date | null => {
  if (!dateString) return null;
  if (dateString instanceof Date) return dateString;
  const date = new Date(dateString);
  return isNaN(date.getTime()) ? null : date;
};

interface TaskListProps {
  tasks: TaskDTO[];
  filter: Filter;
  onToggleTask: (id: number) => void;
  onEditTask: (task: TaskDTO) => void;
  onDeleteTask: (id: number) => void;
  loading?: boolean;
  emptyState?: React.ReactNode;
}

const TaskList: React.FC<TaskListProps> = ({
  tasks,
  filter,
  onToggleTask,
  onEditTask,
  onDeleteTask,
  loading = false,
  emptyState,
}) => {
  const [editingTaskId, setEditingTaskId] = useState<number | null>(null);
  const [editingTask, setEditingTask] = useState<TaskDTO | null>(null);

  const handleEdit = (task: TaskDTO) => {
    setEditingTask(task);
    setEditingTaskId(task.id);
  };

  const handleEditSubmit = async (data: { title: string; priority: Priority; dueDate: Date | string | null }) => {
    if (!editingTask) return;
    
    // Ensure dueDate is properly formatted as ISO string or null
    const dueDate = data.dueDate 
      ? data.dueDate instanceof Date 
        ? data.dueDate.toISOString() 
        : data.dueDate
      : null;
    
    await onEditTask({
      ...editingTask,
      title: data.title,
      priority: data.priority,
      dueDate,
    });
    
    setEditingTaskId(null);
    setEditingTask(null);
  };

  const handleEditCancel = () => {
    setEditingTaskId(null);
    setEditingTask(null);
  };

  const filteredTasks = tasks.filter(task => {
    if (filter === 'active') return !task.completed;
    if (filter === 'completed') return task.completed;
    return true;
  });

  if (loading && tasks.length === 0) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 animate-pulse">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-2"></div>
            <div className="h-3 bg-gray-100 dark:bg-gray-700 rounded w-1/2"></div>
          </div>
        ))}
      </div>
    );
  }

  if (filteredTasks.length === 0) {
    return (
      <div className="text-center py-12">
        {emptyState || (
          <div className="text-gray-500 dark:text-gray-400">
            <svg
              className="mx-auto h-12 w-12 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1}
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
              />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">No tasks found</h3>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {filter === 'all'
                ? 'Get started by creating a new task.'
                : filter === 'active'
                ? 'All tasks are completed!'
                : 'No completed tasks yet.'}
            </p>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {filteredTasks.map((task) =>
        editingTaskId === task.id && editingTask ? (
          <div key={task.id} className="p-4 bg-white dark:bg-gray-800 rounded-lg border border-blue-200 dark:border-blue-800 shadow-sm">
            <TaskForm
              initialData={{
                title: editingTask.title,
                priority: editingTask.priority as Priority,
                dueDate: parseDate(editingTask.dueDate),
              }}
              onSubmit={handleEditSubmit}
              onCancel={handleEditCancel}
              loading={loading}
              submitText="Save Changes"
            />
          </div>
        ) : (
          <TaskItem
            key={task.id}
            task={task}
            onToggle={onToggleTask}
            onEdit={handleEdit}
            onDelete={onDeleteTask}
            loading={loading}
          />
        )
      )}
    </div>
  );
};

export default TaskList;
