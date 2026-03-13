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
      className={`file-item ${isSelected ? 'file-item--selected' : ''}`}
      onClick={onClick}
    >
      <span className={`file-item-icon ${isSelected ? 'opacity-100' : ''}`}>
        {file.isDir ? '📁' : '📄'}
      </span>
      <span className="file-item-name">
        {file.name}
      </span>
    </div>
  );
};
