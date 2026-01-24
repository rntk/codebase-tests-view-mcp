import React from 'react';
import { CodeViewer } from './CodeViewer';
import { MindMap } from '../MindMap/MindMap';
import type { FileContent, MindMapNode } from '../../types';

import { filterItemsByLine } from '../../utils/testUtils';

interface FilePreviewProps {
  file: FileContent | null;
  loading: boolean;
  error: string | null;
  onTestClick?: (testId: string) => void;
  selectedLine?: number | null;
  onLineSelect?: (line: number) => void;
}

export const FilePreview: React.FC<FilePreviewProps> = ({
  file,
  loading,
  error,
  onTestClick,
  selectedLine,
  onLineSelect,
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
  let testRefs = file.metadata?.tests ?? [];

  // Filter testRefs if a line is selected
  testRefs = filterItemsByLine(testRefs, selectedLine, (test) => test.coveredLines);

  const functionMap = new Map<string, typeof testRefs>();
  testRefs.forEach((test) => {
    const functionName = test.functionName.trim();
    const existing = functionMap.get(functionName);
    if (existing) {
      existing.push(test);
    } else {
      functionMap.set(functionName, [test]);
    }
  });

  const mindMapData: MindMapNode = {
    id: file.path,
    label: file.name,
    children: Array.from(functionMap.entries()).map(([functionName, tests]) => ({
      id: `func:${file.path}:${functionName}`,
      label: functionName,
      children: tests.map((test) => ({
        id: `${test.testFile}:${test.testName}`,
        label: test.testName,
        edgeLabel: test.comment,
      })),
    })),
  };

  const hasTests = mindMapData.children && mindMapData.children.length > 0;

  return (
    <div>
      <CodeViewer
        content={file.content}
        filename={file.name}
        testReferences={file.metadata?.tests ?? []}
        coverageDepth={file.coverageDepth}
        onLineClick={onTestClick}
        selectedLine={selectedLine}
        onLineSelect={onLineSelect}
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
