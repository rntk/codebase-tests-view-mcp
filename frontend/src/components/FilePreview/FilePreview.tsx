import React, { useState } from 'react';
import { CodeViewer } from './CodeViewer';
import { MindMap } from '../MindMap/MindMap';
import { TabButton } from '../FileExplorer/TabButton';
import type { FileContent, MindMapNode, Comment, SourceReference } from '../../types';

import { filterItemsByLine } from '../../utils/testUtils';

type PreviewTab = 'code' | 'graph';

interface FilePreviewProps {
  file: FileContent | null;
  loading: boolean;
  error: string | null;
  onTestClick?: (testId: string) => void;
  sourceReferences?: SourceReference[];
  onSourceRefClick?: (sourceFile: string, line: number) => void;
  selectedLine?: number | null;
  onLineSelect?: (line: number) => void;
  onLineDoubleClick?: (line: number) => void;
  onResetLineFilter?: () => void;
  comments?: Comment[];
}

export const FilePreview: React.FC<FilePreviewProps> = ({
  file,
  loading,
  error,
  onTestClick,
  sourceReferences = [],
  onSourceRefClick,
  selectedLine,
  onLineSelect,
  onLineDoubleClick,
  onResetLineFilter,
  comments = [],
}) => {
  const [activeTab, setActiveTab] = useState<PreviewTab>('code');

  if (loading) {
    return <div>Loading file...</div>;
  }

  if (error) {
    return <div className="text-error">Error: {error}</div>;
  }

  if (!file) {
    return <div className="text-secondary">Select a file to view its content</div>;
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

  // Build reverse mind map: test file → source functions → source files
  const hasSourceRefs = sourceReferences.length > 0;
  const reverseMindMapData: MindMapNode | null = hasSourceRefs ? (() => {
    const sourceFileMap = new Map<string, Map<string, SourceReference[]>>();
    sourceReferences.forEach(ref => {
      if (!sourceFileMap.has(ref.sourceFile)) {
        sourceFileMap.set(ref.sourceFile, new Map());
      }
      const funcMap = sourceFileMap.get(ref.sourceFile)!;
      if (!funcMap.has(ref.functionName)) {
        funcMap.set(ref.functionName, []);
      }
      funcMap.get(ref.functionName)!.push(ref);
    });

    return {
      id: file.path,
      label: file.name,
      children: Array.from(sourceFileMap.entries()).map(([sourceFile, funcMap]) => ({
        id: `src:${sourceFile}`,
        label: sourceFile.split('/').pop() ?? sourceFile,
        children: Array.from(funcMap.entries()).map(([funcName, refs]) => ({
          id: `src:${sourceFile}:${funcName}`,
          label: funcName,
          edgeLabel: refs[0]?.comment,
        })),
      })),
    };
  })() : null;

  const hasGraphData = hasTests || hasSourceRefs;

  return (
    <div className="file-preview-container">
      {/* Tab bar */}
      <div className="tab-bar">
        <TabButton
          label="Code"
          isActive={activeTab === 'code'}
          onClick={() => setActiveTab('code')}
        />
        <TabButton
          label="Graph"
          isActive={activeTab === 'graph'}
          onClick={() => setActiveTab('graph')}
        />
      </div>

      {/* Tab content */}
      {activeTab === 'code' && (
        <CodeViewer
          content={file.content}
          filename={file.name}
          testReferences={file.metadata?.tests ?? []}
          coverageDepth={file.coverageDepth}
          comments={comments}
          sourceReferences={sourceReferences}
          onLineClick={onTestClick}
          onSourceRefClick={onSourceRefClick}
          selectedLine={selectedLine}
          onLineSelect={onLineSelect}
          onLineDoubleClick={onLineDoubleClick}
        />
      )}

      {activeTab === 'graph' && (
        <div className="graph-tab-content">
          {hasTests && (
            <div className="mt-lg">
              <div className="file-preview-actions">
                <h3 className="section-title">
                  Test Coverage
                </h3>
                {selectedLine !== null && selectedLine !== undefined && onResetLineFilter && (
                  <button
                    type="button"
                    onClick={onResetLineFilter}
                    className="btn btn-ghost"
                    title="Clear line filter and show all functions"
                  >
                    Show all functions
                  </button>
                )}
              </div>
              <MindMap data={mindMapData} onNodeClick={onTestClick} />
            </div>
          )}

          {hasSourceRefs && reverseMindMapData && (
            <div className="mt-lg">
              <div className="mb-md">
                <h3 className="section-title">
                  Source Coverage
                </h3>
                <p className="text-secondary text-xs m-0">
                  Click highlighted lines to navigate to the source function
                </p>
              </div>
              <MindMap data={reverseMindMapData} onNodeClick={(nodeId) => {
                // nodeId format: "src:sourceFile:funcName" or "src:sourceFile"
                const parts = nodeId.split(':');
                if (parts.length >= 3) {
                  const sourceFile = parts.slice(1, -1).join(':');
                  const ref = sourceReferences.find(r => r.sourceFile === sourceFile);
                  if (ref && onSourceRefClick) {
                    onSourceRefClick(ref.sourceFile, ref.coveredLines.start);
                  }
                } else if (parts.length === 2) {
                  const sourceFile = parts[1];
                  const ref = sourceReferences.find(r => r.sourceFile === sourceFile);
                  if (ref && onSourceRefClick) {
                    onSourceRefClick(ref.sourceFile, ref.coveredLines.start);
                  }
                }
              }} />
            </div>
          )}

          {!hasGraphData && (
            <div className="mt-lg text-secondary text-center">
              No test metadata available for this file
            </div>
          )}
        </div>
      )}
    </div>
  );
};
