import type { Status } from '../types';

const statusConfig: Record<Status, { color: string; bg: string; label: string }> = {
  pending: { color: '#6b7280', bg: '#f3f4f6', label: 'Pending' },
  in_progress: { color: '#2563eb', bg: '#dbeafe', label: 'In Progress' },
  completed: { color: '#16a34a', bg: '#dcfce7', label: 'Completed' },
  failed: { color: '#dc2626', bg: '#fee2e2', label: 'Failed' },
};

interface StatusBadgeProps {
  status: Status;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const config = statusConfig[status];

  return (
    <span
      style={{
        display: 'inline-block',
        padding: '2px 8px',
        borderRadius: '4px',
        fontSize: '12px',
        fontWeight: 500,
        color: config.color,
        backgroundColor: config.bg,
      }}
    >
      {config.label}
    </span>
  );
}
