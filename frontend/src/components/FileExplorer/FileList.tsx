import React from 'react';
import { FileItem } from './FileItem';
import { PathInput } from './PathInput';
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
    <div>
      <h2 style={{ marginTop: 0, marginBottom: '16px', fontSize: '18px' }}>
        Files
      </h2>
      <PathInput path={path} onChange={onPathChange} />

      {loading && <div>Loading...</div>}
      {error && <div style={{ color: 'red' }}>Error: {error}</div>}

      {!loading && !error && (
        <div>
          {files.length === 0 && <div>No files found</div>}
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
