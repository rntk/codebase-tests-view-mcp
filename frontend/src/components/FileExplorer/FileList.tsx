import React, { useState, useMemo } from 'react';
import { FileItem } from './FileItem';
import { PathInput } from './PathInput';
import { Breadcrumbs } from './Breadcrumbs';
import { TabButton } from './TabButton';
import { TestsList } from './TestsList';
import { FunctionsWithTests } from './FunctionsWithTests';
import { ProjectOverview } from './ProjectOverview';
import { MetadataIssuesPanel } from './MetadataIssuesPanel';
import type { FileEntry, TestDetail, OverviewResponse, MetadataIssue } from '../../types';
import { groupTestsByFunction } from '../../utils/testUtils';

type ExplorerTab = 'overview' | 'issues' | 'files';
type FilesSubview = 'files' | 'tests' | 'functions';

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
  // Overview props
  overview: OverviewResponse | null;
  overviewLoading: boolean;
  overviewError: string | null;
  metadataIssues: MetadataIssue[];
  metadataIssuesLoading: boolean;
  metadataIssuesError: string | null;
  onMetadataChanged: () => void;
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
  overview,
  overviewLoading,
  overviewError,
  metadataIssues,
  metadataIssuesLoading,
  metadataIssuesError,
  onMetadataChanged,
}) => {
  const [activeTab, setActiveTab] = useState<ExplorerTab>('overview');
  const [activeFilesSubview, setActiveFilesSubview] = useState<FilesSubview>('files');

  // Group tests by function name for the functions tab (using shared utility)
  const functionsWithTests = useMemo(() => groupTestsByFunction(tests), [tests]);

  return (
    <div className="file-list-container">
      <div className="file-list-header">
        <h2 className="section-title mt-0 mb-md">
          Explorer
        </h2>

        <div className="tab-bar">
          <TabButton
            label="Overview"
            isActive={activeTab === 'overview'}
            onClick={() => setActiveTab('overview')}
            badge={overview?.totalTests}
          />
          <TabButton
            label="Files"
            isActive={activeTab === 'files'}
            onClick={() => setActiveTab('files')}
            badge={files.length}
          />
          <TabButton
            label="Issues"
            isActive={activeTab === 'issues'}
            onClick={() => setActiveTab('issues')}
            badge={metadataIssues.length}
          />
        </div>
      </div>

      {activeTab === 'overview' && (
        <div className="file-list-tab-scroll">
          <ProjectOverview
            overview={overview}
            loading={overviewLoading}
            error={overviewError}
            onTestClick={onTestClick}
            onSourceFileClick={onSourceFileClick}
          />
        </div>
      )}

      {activeTab === 'files' && (
        <div className="file-list-tab-content file-list-tab-content--files">
          <div className="file-list-files-nav">
            <PathInput path={path} onChange={onPathChange} />
            <Breadcrumbs currentPath={path} onPathChange={onPathChange} />
          </div>

          <div className="file-list-subtabs">
            <TabButton
              label="Files"
              isActive={activeFilesSubview === 'files'}
              onClick={() => setActiveFilesSubview('files')}
              badge={files.length}
            />
            <TabButton
              label="Tests"
              isActive={activeFilesSubview === 'tests'}
              onClick={() => setActiveFilesSubview('tests')}
              badge={tests.length}
            />
            <TabButton
              label="Functions"
              isActive={activeFilesSubview === 'functions'}
              onClick={() => setActiveFilesSubview('functions')}
              badge={functionsWithTests.size}
            />
          </div>

          <div className="file-list-files-body">
            {activeFilesSubview === 'files' && (
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

            {activeFilesSubview === 'tests' && (
              <div className="file-list-tab-scroll">
                <TestsList
                  tests={tests}
                  loading={testsLoading}
                  error={testsError}
                  onTestClick={onTestClick}
                  onSourceFileClick={onSourceFileClick}
                />
              </div>
            )}

            {activeFilesSubview === 'functions' && (
              <div className="file-list-tab-scroll">
                <FunctionsWithTests
                  tests={tests}
                  loading={testsLoading}
                  error={testsError}
                  onTestClick={onTestClick}
                  onSourceFileClick={onSourceFileClick}
                />
              </div>
            )}
          </div>
        </div>
      )}

      {activeTab === 'issues' && (
        <div className="file-list-tab-scroll">
          <MetadataIssuesPanel
            issues={metadataIssues}
            loading={metadataIssuesLoading}
            error={metadataIssuesError}
            onMetadataChanged={onMetadataChanged}
            onSourceFileClick={onSourceFileClick}
          />
        </div>
      )}
    </div>
  );
};
