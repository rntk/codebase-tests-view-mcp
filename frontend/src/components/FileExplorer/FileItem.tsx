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
        padding: '8px',
        cursor: 'pointer',
        backgroundColor: isSelected ? '#e3f2fd' : 'transparent',
        borderRadius: '4px',
        marginBottom: '4px',
        display: 'flex',
        alignItems: 'center',
        transition: 'background-color 0.2s',
      }}
      onMouseEnter={(e) => {
        if (!isSelected) {
          e.currentTarget.style.backgroundColor = '#f5f5f5';
        }
      }}
      onMouseLeave={(e) => {
        if (!isSelected) {
          e.currentTarget.style.backgroundColor = 'transparent';
        }
      }}
    >
      <span style={{ marginRight: '8px' }}>
        {file.isDir ? '📁' : '📄'}
      </span>
      <span style={{ fontFamily: 'monospace', fontSize: '14px' }}>
        {file.name}
      </span>
    </div>
  );
};
