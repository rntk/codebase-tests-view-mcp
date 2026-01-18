import React from 'react';
import { CodeViewer } from './CodeViewer';
import { MindMap } from '../MindMap/MindMap';
import type { FileContent, MindMapNode } from '../../types';

interface FilePreviewProps {
  file: FileContent | null;
  loading: boolean;
  error: string | null;
  onTestClick?: (testId: string) => void;
}

export const FilePreview: React.FC<FilePreviewProps> = ({
  file,
  loading,
  error,
  onTestClick,
}) => {
  if (loading) {
    return <div>Loading file...</div>;
  }

  if (error) {
    return <div style={{ color: 'red' }}>Error: {error}</div>;
  }

  if (!file) {
    return <div style={{ color: '#666' }}>Select a file to view its content</div>;
  }

  // Build mind map data from file metadata
  const testRefs = file.metadata?.tests ?? [];
  const mindMapData: MindMapNode = {
    id: file.path,
    label: file.name,
    children: testRefs.map((test) => ({
      id: `${test.testFile}:${test.testName}`,
      label: test.testName,
      edgeLabel: test.comment,
    })),
  };

  const hasTests = mindMapData.children && mindMapData.children.length > 0;

  return (
    <div>
      <CodeViewer
        content={file.content}
        filename={file.name}
        testReferences={testRefs}
        onLineClick={onTestClick}
      />

      {hasTests && (
        <div style={{ marginTop: '24px' }}>
          <h3 style={{ marginBottom: '16px', fontSize: '16px' }}>
            Test Coverage
          </h3>
          <MindMap data={mindMapData} onNodeClick={onTestClick} />
        </div>
      )}

      {!hasTests && (
        <div style={{ marginTop: '24px', color: '#666', textAlign: 'center' }}>
          No test metadata available for this file
        </div>
      )}
    </div>
  );
};
