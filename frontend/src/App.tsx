import { useState } from 'react';
import { ThreePanel } from './components/Layout/ThreePanel';
import { FileList } from './components/FileExplorer/FileList';
import { FilePreview } from './components/FilePreview/FilePreview';
import { TestPanel } from './components/TestPanel/TestPanel';
import { useFiles } from './hooks/useFiles';
import { useFileContent } from './hooks/useFileContent';
import { useTests } from './hooks/useTests';
import type { FileEntry } from './types';

function App() {
  const [currentPath, setCurrentPath] = useState('.');
  const [selectedFilePath, setSelectedFilePath] = useState<string | null>(null);
  const [highlightedTestIds, setHighlightedTestIds] = useState<Set<string>>(new Set());

  // Load files for current directory
  const { files, loading: filesLoading, error: filesError } = useFiles(currentPath);

  // Load selected file content
  const { file, loading: fileLoading, error: fileError } = useFileContent(selectedFilePath);

  // Load tests for selected file
  const { tests, loading: testsLoading, error: testsError } = useTests(selectedFilePath);

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
        <TestPanel
          tests={tests}
          loading={testsLoading}
          error={testsError}
          highlightedTestIds={highlightedTestIds}
        />
      }
    />
  );
}

export default App;
