import { useState } from 'react';
import { ThreePanel } from './components/Layout/ThreePanel';
import { FileList } from './components/FileExplorer/FileList';
import { FilePreview } from './components/FilePreview/FilePreview';
import { TestPanel } from './components/TestPanel/TestPanel';
import { SuggestionsPanel } from './components/SuggestionsPanel/SuggestionsPanel';
import { useFiles } from './hooks/useFiles';
import { useFileContent } from './hooks/useFileContent';
import { useTests } from './hooks/useTests';
import { useSuggestions } from './hooks/useSuggestions';
import type { FileEntry } from './types';

type RightPanelTab = 'tests' | 'suggestions';

function App() {
  const [currentPath, setCurrentPath] = useState('.');
  const [selectedFilePath, setSelectedFilePath] = useState<string | null>(null);
  const [highlightedTestIds, setHighlightedTestIds] = useState<Set<string>>(new Set());
  const [activeRightTab, setActiveRightTab] = useState<RightPanelTab>('tests');

  // Load files for current directory
  const { files, loading: filesLoading, error: filesError } = useFiles(currentPath);

  // Load selected file content
  const { file, loading: fileLoading, error: fileError } = useFileContent(selectedFilePath);

  // Load tests for selected file
  const { tests, loading: testsLoading, error: testsError } = useTests(selectedFilePath);

  // Load suggestions for selected file
  const { suggestions, loading: suggestionsLoading, error: suggestionsError } = useSuggestions(selectedFilePath);

  // Handle file/directory click
  const handleFileClick = (fileEntry: FileEntry) => {
    if (fileEntry.isDir) {
      // Navigate to directory
      setCurrentPath(fileEntry.path);
      setSelectedFilePath(null);
    } else {
      // Select file
      setSelectedFilePath(fileEntry.path);
    }
    // Clear highlighted tests when changing files
    setHighlightedTestIds(new Set());
  };

  // Handle path change
  const handlePathChange = (newPath: string) => {
    setCurrentPath(newPath);
    setSelectedFilePath(null);
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

  const tabButtonStyle = (isActive: boolean): React.CSSProperties => ({
    padding: '8px 16px',
    border: 'none',
    borderBottom: isActive ? '2px solid var(--primary)' : '2px solid transparent',
    backgroundColor: 'transparent',
    color: isActive ? 'var(--primary)' : 'var(--text-secondary)',
    cursor: 'pointer',
    fontSize: '13px',
    fontWeight: isActive ? '600' : '500',
    transition: 'all 0.15s ease',
  });

  return (
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
        />
      }
      center={
        <FilePreview
          file={file}
          loading={fileLoading}
          error={fileError}
          onTestClick={handleTestClick}
        />
      }
      right={
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
          {/* Tab buttons */}
          <div style={{
            display: 'flex',
            borderBottom: '1px solid var(--border-color)',
            backgroundColor: 'var(--bg-secondary)',
          }}>
            <button
              type="button"
              style={tabButtonStyle(activeRightTab === 'tests')}
              onClick={() => setActiveRightTab('tests')}
            >
              Tests
              {tests.length > 0 && (
                <span style={{
                  marginLeft: '6px',
                  padding: '1px 6px',
                  backgroundColor: activeRightTab === 'tests' ? 'var(--primary)' : 'var(--bg-tertiary)',
                  color: activeRightTab === 'tests' ? 'white' : 'var(--text-tertiary)',
                  borderRadius: '10px',
                  fontSize: '11px',
                }}>
                  {tests.length}
                </span>
              )}
            </button>
            <button
              type="button"
              style={tabButtonStyle(activeRightTab === 'suggestions')}
              onClick={() => setActiveRightTab('suggestions')}
            >
              Suggestions
              {suggestions.length > 0 && (
                <span style={{
                  marginLeft: '6px',
                  padding: '1px 6px',
                  backgroundColor: activeRightTab === 'suggestions' ? 'var(--primary)' : 'var(--bg-tertiary)',
                  color: activeRightTab === 'suggestions' ? 'white' : 'var(--text-tertiary)',
                  borderRadius: '10px',
                  fontSize: '11px',
                }}>
                  {suggestions.length}
                </span>
              )}
            </button>
          </div>

          {/* Tab content */}
          <div style={{ flex: 1, overflow: 'auto', padding: 'var(--space-md)' }}>
            {activeRightTab === 'tests' && (
              <TestPanel
                tests={tests}
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
          </div>
        </div>
      }
    />
  );
}

export default App;
