import React from 'react';
import type { FileEntry } from '../../types';

interface FileItemProps {
  file: FileEntry;
  onClick: () => void;
  isSelected: boolean;
}

export const FileItem: React.FC<FileItemProps> = ({ file, onClick, isSelected }) => {
  return (
    <div
      onClick={onClick}
      style={{
        padding: 'var(--space-sm) var(--space-md)',
        cursor: 'pointer',
        backgroundColor: isSelected ? 'var(--accent-soft)' : 'transparent',
        color: isSelected ? 'var(--accent-secondary)' : 'var(--text-primary)',
        borderRadius: 'var(--radius-md)',
        marginBottom: '2px',
        display: 'flex',
        alignItems: 'center',
        transition: 'all 0.2s',
        fontWeight: isSelected ? '600' : '400',
      }}
      onMouseEnter={(e) => {
        if (!isSelected) {
          e.currentTarget.style.backgroundColor = 'var(--bg-tertiary)';
        }
      }}
      onMouseLeave={(e) => {
        if (!isSelected) {
          e.currentTarget.style.backgroundColor = 'transparent';
        }
      }}
    >
      <span style={{
        marginRight: 'var(--space-sm)',
        fontSize: '16px',
        opacity: isSelected ? 1 : 0.7
      }}>
        {file.isDir ? '📁' : '📄'}
      </span>
      <span style={{
        fontFamily: 'var(--font-mono)',
        fontSize: '13px',
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis'
      }}>
        {file.name}
      </span>
    </div>
  );
};
