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
      onClick={onClick}
      style={{
        padding: '8px 12px',
        border: 'none',
        borderBottom: isActive ? '2px solid var(--accent-primary)' : '2px solid transparent',
        backgroundColor: 'transparent',
        color: isActive ? 'var(--accent-primary)' : 'var(--text-secondary)',
        cursor: 'pointer',
        fontSize: '13px',
        fontWeight: isActive ? '600' : '500',
        transition: 'all 0.15s ease',
        display: 'flex',
        alignItems: 'center',
        gap: '6px',
      }}
    >
      {label}
      {badge !== undefined && badge > 0 && (
        <span
          style={{
            padding: '1px 6px',
            backgroundColor: isActive ? 'var(--accent-primary)' : 'var(--bg-tertiary)',
            color: isActive ? 'white' : 'var(--text-tertiary)',
            borderRadius: '10px',
            fontSize: '11px',
            minWidth: '18px',
            textAlign: 'center',
          }}
        >
          {badge}
        </span>
      )}
    </button>
  );
};
