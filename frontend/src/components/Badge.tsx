import React from 'react';

export type BadgeVariant =
  | 'open'
  | 'in_progress'
  | 'resolved'
  | 'closed'
  | 'low'
  | 'medium'
  | 'high'
  | 'urgent'
  | 'active'
  | 'in_service'
  | 'maintenance'
  | 'decommissioned'
  | 'retired'
  | 'available'
  | 'admin'
  | 'technician'
  | 'employee';

interface BadgeProps {
  variant?: BadgeVariant | string;
  children: React.ReactNode;
}

export const Badge: React.FC<BadgeProps> = ({ variant = 'open', children }) => {
  const getStyles = () => {
    switch (variant) {
      case 'open':
      case 'active':
      case 'in_service':
      case 'available':
        return 'bg-emerald-50 text-emerald-700 border-emerald-200';
      case 'in_progress':
      case 'medium':
        return 'bg-amber-50 text-amber-700 border-amber-200';
      case 'resolved':
      case 'low':
        return 'bg-blue-50 text-blue-700 border-blue-200';
      case 'closed':
      case 'decommissioned':
      case 'retired':
        return 'bg-slate-100 text-slate-600 border-slate-200';
      case 'high':
      case 'maintenance':
        return 'bg-orange-50 text-orange-700 border-orange-200';
      case 'urgent':
      case 'admin':
        return 'bg-rose-50 text-rose-700 border-rose-200';
      case 'technician':
        return 'bg-purple-50 text-purple-700 border-purple-200';
      case 'employee':
      default:
        return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStyles()}`}
    >
      {children}
    </span>
  );
};
