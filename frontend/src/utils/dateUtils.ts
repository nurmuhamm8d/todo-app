// Helper function to convert date string to Date object
export const parseDate = (dateString: string | Date | null | undefined): Date | null => {
  if (!dateString) return null;
  if (dateString instanceof Date) return dateString;
  const date = new Date(dateString);
  return isNaN(date.getTime()) ? null : date;
};

export const formatDate = (dateInput?: string | Date | null): string => {
  const date = parseDate(dateInput);
  if (!date) return '';
  
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
};

export const formatRelativeTime = (dateInput: string | Date): string => {
  const date = parseDate(dateInput);
  if (!date) return '';
  
  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));
  
  if (diffInMs < 0) {
    return 'In the future';
  } else if (diffInDays === 0) {
    return 'Today';
  } else if (diffInDays === 1) {
    return 'Yesterday';
  } else if (diffInDays < 7) {
    return `${diffInDays} days ago`;
  } else if (diffInDays < 30) {
    const weeks = Math.floor(diffInDays / 7);
    return `${weeks} ${weeks === 1 ? 'week' : 'weeks'} ago`;
  } else {
    return formatDate(date);
  }
};

export const isDueSoon = (dueDateInput?: string | Date | null, days = 1): boolean => {
  const dueDate = parseDate(dueDateInput);
  if (!dueDate) return false;
  
  const now = new Date();
  const diffInMs = dueDate.getTime() - now.getTime();
  const diffInDays = diffInMs / (1000 * 60 * 60 * 24);
  
  return diffInDays > 0 && diffInDays <= days;
};

export const isOverdue = (dueDateInput?: string | Date | null): boolean => {
  const dueDate = parseDate(dueDateInput);
  if (!dueDate) return false;
  
  return new Date() > dueDate;
};

export const getPriorityColor = (priority: string): string => {
  switch (priority.toLowerCase()) {
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
