import React from 'react';
import { FileItem } from './FileItem';
import { PathInput } from './PathInput';
import { Breadcrumbs } from './Breadcrumbs';
import type { FileEntry } from '../../types';

interface FileListProps {
  path: string;
  files: FileEntry[];
  selectedPath: string | null;
  loading: boolean;
  error: string | null;
  onPathChange: (path: string) => void;
  onFileClick: (file: FileEntry) => void;
}

export const FileList: React.FC<FileListProps> = ({
  path,
  files,
  selectedPath,
  loading,
  error,
  onPathChange,
  onFileClick,
}) => {
  return (
    <div className="file-list-container">
      <h2 style={{
        marginTop: 0,
        marginBottom: 'var(--space-md)',
        fontSize: '16px',
        fontWeight: '600',
        color: 'var(--text-primary)'
      }}>
        Explorer
      </h2>

      <Breadcrumbs currentPath={path} onPathChange={onPathChange} />

      <div style={{ marginBottom: 'var(--space-md)' }}>
        <PathInput path={path} onChange={onPathChange} />
      </div>

      {loading && (
        <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)' }}>
          Loading files...
        </div>
      )}

      {error && (
        <div style={{
          padding: 'var(--space-md)',
          color: 'var(--error)',
          backgroundColor: '#fef2f2',
          borderRadius: 'var(--radius-md)',
          fontSize: '13px'
        }}>
          Error: {error}
        </div>
      )}

      {!loading && !error && (
        <div className="file-items">
          {files.length === 0 && (
            <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)', textAlign: 'center' }}>
              No files found in this directory
            </div>
          )}
          {files.map((file) => (
            <FileItem
              key={file.path}
              file={file}
              onClick={() => onFileClick(file)}
              isSelected={selectedPath === file.path}
            />
          ))}
        </div>
      )}
    </div>
  );
};
