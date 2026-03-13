import React from 'react';

interface TabButtonProps {
  label: string;
  isActive: boolean;
  onClick: () => void;
  badge?: number;
}

export const TabButton: React.FC<TabButtonProps> = ({
  label,
  isActive,
  onClick,
  badge,
}) => {
  return (
    <button
      type="button"
      className={`tab-button ${isActive ? 'tab-button--active' : ''}`}
      onClick={onClick}
    >
      {label}
      {badge !== undefined && badge > 0 && (
        <span className={`tab-badge ${isActive ? 'tab-badge--active' : ''}`}>
          {badge}
        </span>
      )}
    </button>
  );
};
