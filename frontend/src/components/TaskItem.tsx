import React from 'react';
import { formatDate, isDueSoon, isOverdue } from '../utils/dateUtils';
import { Priority, TaskDTO } from '../types';

interface TaskItemProps {
  task: TaskDTO;
  onToggle: (id: number) => void;
  onEdit: (task: TaskDTO) => void;
  onDelete: (id: number) => void;
  loading?: boolean;
}

const getPriorityLabel = (priority: string) => {
  switch (priority.toLowerCase()) {
    case 'high':
      return 'High';
    case 'medium':
      return 'Medium';
    case 'low':
      return 'Low';
    default:
      return priority;
  }
};

const getPriorityColor = (priority: Priority) => {
  switch (priority) {
    case 'high':
      return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
    case 'medium':
      return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
    case 'low':
      return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
  }
};

export const TaskItem: React.FC<TaskItemProps> = ({
  task,
  onToggle,
  onEdit,
  onDelete,
  loading = false,
}) => {
  const handleToggle = () => {
    if (!loading) {
      onToggle(task.id);
    }
  };

  const handleEdit = () => {
    if (!loading) {
      onEdit(task);
    }
  };

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!loading && window.confirm('Are you sure you want to delete this task?')) {
      onDelete(task.id);
    }
  };

  const dueDate = task.dueDate ? (
    <span className={`text-sm ${
      task.completed 
        ? 'text-gray-500 dark:text-gray-400' 
        : isOverdue(task.dueDate) 
          ? 'text-red-600 dark:text-red-400' 
          : isDueSoon(task.dueDate) 
            ? 'text-yellow-600 dark:text-yellow-400' 
            : 'text-gray-500 dark:text-gray-400'
    }`}>
      {formatDate(task.dueDate)}
    </span>
  ) : null;

  return (
    <div
      className={`group relative p-4 border rounded-lg transition-all duration-200 ${
        task.completed
          ? 'bg-gray-50 dark:bg-gray-800 border-gray-200 dark:border-gray-700 opacity-75'
          : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 hover:shadow-md hover:border-blue-200 dark:hover:border-blue-800'
      }`}
    >
      <div className="flex items-start">
        <div className="flex items-center h-5 mt-0.5">
          <input
            type="checkbox"
            checked={task.completed}
            onChange={handleToggle}
            className={`h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 ${
              loading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
            }`}
            disabled={loading}
            aria-label={task.completed ? 'Mark as incomplete' : 'Mark as complete'}
          />
        </div>

        <div className="ml-3 flex-1 min-w-0">
          <div className="flex justify-between items-start">
            <div className="flex-1 min-w-0">
              <p
                className={`text-sm font-medium ${
                  task.completed
                    ? 'text-gray-500 dark:text-gray-400 line-through'
                    : 'text-gray-900 dark:text-white'
                }`}
              >
                {task.title}
              </p>
              
              {dueDate && (
                <div className="mt-1">
                  {dueDate}
                </div>
              )}
            </div>

            <div className="flex items-center space-x-2 ml-2">
              <span
                className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityColor(
                  task.priority as Priority
                )}`}
              >
                {getPriorityLabel(task.priority)}
              </span>
            </div>
          </div>
        </div>

        <div className="ml-4 flex-shrink-0 flex space-x-2">
          <button
            type="button"
            onClick={handleEdit}
            disabled={loading}
            className="p-1 text-gray-400 hover:text-blue-500 dark:hover:text-blue-400 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
            aria-label="Edit task"
          >
            <svg
              className="h-5 w-5"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
              aria-hidden="true"
            >
              <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
            </svg>
          </button>
          
          <button
            type="button"
            onClick={handleDelete}
            disabled={loading}
            className="p-1 text-gray-400 hover:text-red-500 dark:hover:text-red-400 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
            aria-label="Delete task"
          >
            <svg
              className="h-5 w-5"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                fillRule="evenodd"
                d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
};
