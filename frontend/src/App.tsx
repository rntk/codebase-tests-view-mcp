import { useState, useEffect, useCallback } from 'react';
import { ThreePanel } from './components/Layout/ThreePanel';
import { FileList } from './components/FileExplorer/FileList';
import { FilePreview } from './components/FilePreview/FilePreview';
import { TestPanel } from './components/TestPanel/TestPanel';
import { SuggestionsPanel } from './components/SuggestionsPanel/SuggestionsPanel';
import { CommentPanel } from './components/CommentPanel';
import { ExportModal } from './components/ExportModal';
import { GlobalSearchPalette } from './components/GlobalSearchPalette';
import { useFiles } from './hooks/useFiles';
import { useFileContent } from './hooks/useFileContent';
import { useTests } from './hooks/useTests';
import { useDirectoryTests } from './hooks/useDirectoryTests';
import { useSuggestions } from './hooks/useSuggestions';
import { useComments } from './hooks/useComments';
import { useExport } from './hooks/useExport';
import { useSources } from './hooks/useSources';
import { useProjectOverview } from './hooks/useProjectOverview';
import { useMetadataIssues } from './hooks/useMetadataIssues';
import type { FileEntry, TestDetail, SearchResult } from './types';
import { filterItemsByLine } from './utils/testUtils';

type RightPanelTab = 'tests' | 'suggestions' | 'comments';

function App() {
  // Initialize state from URL query parameters
  const [currentPath, setCurrentPath] = useState(() => {
    const params = new URLSearchParams(window.location.search);
    return params.get('path') || '.';
  });
  const [selectedFilePath, setSelectedFilePath] = useState<string | null>(() => {
    const params = new URLSearchParams(window.location.search);
    return params.get('file');
  });
  const [selectedLine, setSelectedLine] = useState<number | null>(() => {
    const params = new URLSearchParams(window.location.search);
    const line = params.get('line');
    return line ? parseInt(line, 10) : null;
  });
  const [highlightedTestIds, setHighlightedTestIds] = useState<Set<string>>(new Set());
  const [activeRightTab, setActiveRightTab] = useState<RightPanelTab>('tests');
  const [isExportModalOpen, setIsExportModalOpen] = useState(false);
  const [isSearchPaletteOpen, setIsSearchPaletteOpen] = useState(false);

  // Handle browser back/forward buttons
  useEffect(() => {
    const handlePopState = () => {
      const params = new URLSearchParams(window.location.search);
      const urlPath = params.get('path') || '.';
      const urlFile = params.get('file');
      const urlLine = params.get('line');

      setCurrentPath(urlPath);
      setSelectedFilePath(urlFile);
      setSelectedLine(urlLine ? parseInt(urlLine, 10) : null);
    };

    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  }, []);

  // Global keyboard shortcut for search palette (Cmd+P / Ctrl+P)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Check for Cmd+P or Ctrl+P
      if ((e.metaKey || e.ctrlKey) && e.key === 'p') {
        e.preventDefault();
        setIsSearchPaletteOpen(prev => !prev);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Update URL when state changes
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const urlPath = params.get('path') || '.';
    const urlFile = params.get('file');
    const urlLine = params.get('line');
    const stateLine = selectedLine?.toString() || null;

    // Only push if URL doesn't match state to avoid infinite loops and unnecessary history entries
    if (urlPath !== currentPath || urlFile !== selectedFilePath || urlLine !== stateLine) {
      const newParams = new URLSearchParams();
      newParams.set('path', currentPath);
      if (selectedFilePath) {
        newParams.set('file', selectedFilePath);
      }
      if (selectedLine) {
        newParams.set('line', selectedLine.toString());
      }

      const newUrl = `${window.location.pathname}?${newParams.toString()}`;
      window.history.pushState(null, '', newUrl);
    }
  }, [currentPath, selectedFilePath, selectedLine]);

  // Load files for current directory
  const { files, loading: filesLoading, error: filesError } = useFiles(currentPath);

  // Load tests for all files in the current directory (for the Explorer tabs)
  const { tests: directoryTests, loading: directoryTestsLoading, error: directoryTestsError } = useDirectoryTests(currentPath, files);

  // Load project overview (global summary of all tests and functions)
  const {
    overview,
    loading: overviewLoading,
    error: overviewError,
    refresh: refreshOverview,
  } = useProjectOverview();

  const {
    issues: metadataIssues,
    loading: metadataIssuesLoading,
    error: metadataIssuesError,
    refresh: refreshMetadataIssues,
  } = useMetadataIssues();

  // Load selected file content
  const { file, loading: fileLoading, error: fileError } = useFileContent(selectedFilePath);

  // Load tests for selected file
  const { tests, loading: testsLoading, error: testsError } = useTests(selectedFilePath);

  // Load suggestions for selected file
  const { suggestions, loading: suggestionsLoading, error: suggestionsError } = useSuggestions(selectedFilePath);

  // Load comments for selected file
  const {
    comments,
    loading: commentsLoading,
    error: commentsError,
    addComment,
    editComment,
    removeComment,
    toggleResolved,
  } = useComments(selectedFilePath);

  // Load source references for test files (reverse lookup)
  const { sources: sourceReferences } = useSources(selectedFilePath);

  // Export functionality
  const { exportData, loading: exportLoading, performExport, clearExport } = useExport();

  // Handle file/directory click
  const handleFileClick = (fileEntry: FileEntry) => {
    if (fileEntry.isDir) {
      // Navigate to directory
      setCurrentPath(fileEntry.path);
      setSelectedFilePath(null);
      setSelectedLine(null);
    } else {
      // Select file
      setSelectedFilePath(fileEntry.path);
      setSelectedLine(null);
    }
    // Clear highlighted tests when changing files
    setHighlightedTestIds(new Set());
  };

  // Handle path change
  const handlePathChange = (newPath: string) => {
    setCurrentPath(newPath);
    setSelectedFilePath(null);
    setSelectedLine(null);
    setHighlightedTestIds(new Set());
  };

  // Handle test node click in mind map or code line click
  const handleTestClick = (testId: string) => {
    setHighlightedTestIds(prev => {
      const newSet = new Set(prev);
      if (newSet.has(testId)) {
        newSet.delete(testId);
      } else {
        newSet.add(testId);
      }
      return newSet;
    });
  };

  // Handle test click from Explorer tabs
  const handleExplorerTestClick = useCallback((test: TestDetail) => {
    // Navigate to the test file
    setSelectedFilePath(test.testFile);
    setSelectedLine(test.lineRange?.start ?? null);
    setHighlightedTestIds(new Set());
  }, []);

  // Handle line selection for filtering
  const handleLineSelect = (lineNum: number) => {
    setSelectedLine(prev => {
      if (prev === lineNum) {
        // Deselect if same line clicked
        return null;
      }
      return lineNum;
    });
  };

  // Handle line double click for adding comments
  const handleLineDoubleClick = (lineNum: number) => {
    setSelectedLine(lineNum);
    setActiveRightTab('comments');
  };

  const handleResetLineFilter = () => {
    setSelectedLine(null);
  };

  // Handle click on a source reference line in a test file: navigate to source file at the covered line
  const handleSourceRefClick = (sourceFile: string, line: number) => {
    setSelectedFilePath(sourceFile);
    setSelectedLine(line);
    setHighlightedTestIds(new Set());
  };

  // Handle search result selection
  const handleSearchResultSelect = useCallback((result: SearchResult) => {
    setSelectedFilePath(result.path);
    setSelectedLine(result.line > 0 ? result.line : null);
    setHighlightedTestIds(new Set());
  }, []);

  // Handle export for AI
  const handleExportForAI = async () => {
    if (!selectedFilePath) return;
    
    setIsExportModalOpen(true);
    try {
      await performExport(selectedFilePath, {
        includeTests: true,
        includeSuggestions: true,
        contextLines: 5,
      });
    } catch (err) {
      // Error is handled by the hook
    }
  };

  const handleCloseExportModal = () => {
    setIsExportModalOpen(false);
    clearExport();
  };

  // Filter tests based on selected line
  const filteredTests = selectedLine && file?.metadata?.tests
    ? filterItemsByLine(tests, selectedLine, (test) => {
      // Find matching test reference in file metadata to get covered lines
      const testRef = file.metadata?.tests?.find(
        ref => ref.testFile === test.testFile && ref.testName === test.testName
      );
      return testRef?.coveredLines;
    })
    : tests;

  const unresolvedCount = comments.filter(c => !c.resolved).length;
  const handleMetadataChanged = useCallback(() => {
    refreshOverview();
    refreshMetadataIssues();
  }, [refreshMetadataIssues, refreshOverview]);

  return (
    <>
      <ThreePanel
        left={
          <FileList
            path={currentPath}
            files={files}
            selectedPath={selectedFilePath}
            loading={filesLoading}
            error={filesError}
            onPathChange={handlePathChange}
            onFileClick={handleFileClick}
            tests={directoryTests}
            testsLoading={directoryTestsLoading}
            testsError={directoryTestsError}
            onTestClick={handleExplorerTestClick}
            onSourceFileClick={handleSourceRefClick}
            overview={overview}
            overviewLoading={overviewLoading}
            overviewError={overviewError}
            metadataIssues={metadataIssues}
            metadataIssuesLoading={metadataIssuesLoading}
            metadataIssuesError={metadataIssuesError}
            onMetadataChanged={handleMetadataChanged}
          />
        }
        center={
          <FilePreview
            file={file}
            loading={fileLoading}
            error={fileError}
            onTestClick={handleTestClick}
            sourceReferences={sourceReferences}
            onSourceRefClick={handleSourceRefClick}
            selectedLine={selectedLine}
            onLineSelect={handleLineSelect}
            onLineDoubleClick={handleLineDoubleClick}
            onResetLineFilter={handleResetLineFilter}
            comments={comments}
          />
        }
        right={
          <div className="right-panel-wrapper">
            {/* Tab buttons */}
            <div className="tab-bar">
              <button
                type="button"
                className={`tab-button ${activeRightTab === 'tests' ? 'tab-button--active' : ''}`}
                onClick={() => setActiveRightTab('tests')}
              >
                Tests
                {filteredTests.length > 0 && (
                  <span className={`tab-badge ${activeRightTab === 'tests' ? 'tab-badge--active' : ''}`}>
                    {filteredTests.length}
                  </span>
                )}
              </button>
              <button
                type="button"
                className={`tab-button ${activeRightTab === 'suggestions' ? 'tab-button--active' : ''}`}
                onClick={() => setActiveRightTab('suggestions')}
              >
                Suggestions
                {suggestions.length > 0 && (
                  <span className={`tab-badge ${activeRightTab === 'suggestions' ? 'tab-badge--active' : ''}`}>
                    {suggestions.length}
                  </span>
                )}
              </button>
              <button
                type="button"
                className={`tab-button ${activeRightTab === 'comments' ? 'tab-button--active' : ''}`}
                onClick={() => setActiveRightTab('comments')}
              >
                Comments
                {unresolvedCount > 0 && (
                  <span className={`tab-badge ${activeRightTab === 'comments' ? 'tab-badge--warning' : ''}`}>
                    {unresolvedCount}
                  </span>
                )}
              </button>
            </div>

            {/* Tab content */}
            <div className="tab-content">
              {activeRightTab === 'tests' && (
                <TestPanel
                  tests={filteredTests}
                  loading={testsLoading}
                  error={testsError}
                  highlightedTestIds={highlightedTestIds}
                />
              )}
              {activeRightTab === 'suggestions' && (
                <SuggestionsPanel
                  suggestions={suggestions}
                  loading={suggestionsLoading}
                  error={suggestionsError}
                />
              )}
              {activeRightTab === 'comments' && (
                <CommentPanel
                  comments={comments}
                  loading={commentsLoading}
                  error={commentsError}
                  selectedLine={selectedLine}
                  filePath={selectedFilePath}
                  fileContent={file?.content ?? null}
                  onAddComment={addComment}
                  onUpdateComment={editComment}
                  onDeleteComment={removeComment}
                  onToggleResolved={toggleResolved}
                  onExportForAI={handleExportForAI}
                />
              )}
            </div>
          </div>
        }
      />

      <ExportModal
        isOpen={isExportModalOpen}
        onClose={handleCloseExportModal}
        exportData={exportData}
        loading={exportLoading}
      />

      <GlobalSearchPalette
        isOpen={isSearchPaletteOpen}
        onClose={() => setIsSearchPaletteOpen(false)}
        onResultSelect={handleSearchResultSelect}
      />
    </>
  );
}

export default App;
