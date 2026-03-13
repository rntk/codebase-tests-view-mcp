import React, { useState, useMemo } from 'react';
import { FileItem } from './FileItem';
import { PathInput } from './PathInput';
import { Breadcrumbs } from './Breadcrumbs';
import { TabButton } from './TabButton';
import { TestsList } from './TestsList';
import { FunctionsWithTests } from './FunctionsWithTests';
import type { FileEntry, TestDetail } from '../../types';
import { groupTestsByFunction } from '../../utils/testUtils';

type ExplorerTab = 'files' | 'tests' | 'functions';

interface FileListProps {
  path: string;
  files: FileEntry[];
  selectedPath: string | null;
  loading: boolean;
  error: string | null;
  onPathChange: (path: string) => void;
  onFileClick: (file: FileEntry) => void;
  tests: TestDetail[];
  testsLoading: boolean;
  testsError: string | null;
  onTestClick: (test: TestDetail) => void;
  onSourceFileClick: (sourceFile: string, line: number) => void;
}

export const FileList: React.FC<FileListProps> = ({
  path,
  files,
  selectedPath,
  loading,
  error,
  onPathChange,
  onFileClick,
  tests,
  testsLoading,
  testsError,
  onTestClick,
  onSourceFileClick,
}) => {
  const [activeTab, setActiveTab] = useState<ExplorerTab>('files');

  // Group tests by function name for the functions tab (using shared utility)
  const functionsWithTests = useMemo(() => groupTestsByFunction(tests), [tests]);

  return (
    <div className="file-list-container">
      <h2 className="section-title mt-0 mb-md">
        Explorer
      </h2>

      <Breadcrumbs currentPath={path} onPathChange={onPathChange} />

      <div className="mb-md">
        <PathInput path={path} onChange={onPathChange} />
      </div>

      {/* Tab buttons */}
      <div className="tab-bar">
        <TabButton
          label="Files"
          isActive={activeTab === 'files'}
          onClick={() => setActiveTab('files')}
          badge={files.length}
        />
        <TabButton
          label="Tests"
          isActive={activeTab === 'tests'}
          onClick={() => setActiveTab('tests')}
          badge={tests.length}
        />
        <TabButton
          label="Functions"
          isActive={activeTab === 'functions'}
          onClick={() => setActiveTab('functions')}
          badge={functionsWithTests.size}
        />
      </div>

      {/* Tab content */}
      {activeTab === 'files' && (
        <>
          {loading && (
            <div className="loading-state">
              Loading files...
            </div>
          )}

          {error && (
            <div className="error-state">
              Error: {error}
            </div>
          )}

          {!loading && !error && (
            <div className="file-items">
              {files.length === 0 && (
                <div className="empty-state">
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
        </>
      )}

      {activeTab === 'tests' && (
        <TestsList
          tests={tests}
          loading={testsLoading}
          error={testsError}
          onTestClick={onTestClick}
          onSourceFileClick={onSourceFileClick}
        />
      )}

      {activeTab === 'functions' && (
        <FunctionsWithTests
          tests={tests}
          loading={testsLoading}
          error={testsError}
          onTestClick={onTestClick}
          onSourceFileClick={onSourceFileClick}
        />
      )}
    </div>
  );
};
