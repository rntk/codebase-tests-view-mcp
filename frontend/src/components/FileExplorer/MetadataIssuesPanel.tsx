import React, { useState } from 'react';
import {
  deleteSourcePath,
  deleteTestPath,
  updateSourcePath,
  updateTestPath,
} from '../../api/client';
import type { MetadataIssue, MetadataTestIssue } from '../../types';

interface MetadataIssuesPanelProps {
  issues: MetadataIssue[];
  loading: boolean;
  error: string | null;
  onMetadataChanged: () => void;
  onSourceFileClick: (sourceFile: string, line: number) => void;
}

function testIssueKey(sourceFile: string, testIssue: MetadataTestIssue): string {
  return `${sourceFile}::${testIssue.testFile}::${testIssue.testName}`;
}

export const MetadataIssuesPanel: React.FC<MetadataIssuesPanelProps> = ({
  issues,
  loading,
  error,
  onMetadataChanged,
  onSourceFileClick,
}) => {
  const [sourceDrafts, setSourceDrafts] = useState<Record<string, string>>({});
  const [testDrafts, setTestDrafts] = useState<Record<string, string>>({});
  const [busyKey, setBusyKey] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  if (loading) {
    return <div className="loading-state">Loading metadata issues...</div>;
  }

  if (error) {
    return <div className="error-state">Error: {error}</div>;
  }

  if (issues.length === 0) {
    return (
      <div className="empty-state">
        No metadata path issues found.
      </div>
    );
  }

  const handleSourceSave = async (sourceFile: string) => {
    const nextPath = (sourceDrafts[sourceFile] ?? sourceFile).trim();
    if (!nextPath) {
      setActionError('Replacement source path is required');
      return;
    }

    setBusyKey(`source-save:${sourceFile}`);
    setActionError(null);
    try {
      await updateSourcePath(sourceFile, nextPath);
      setSourceDrafts((current) => ({ ...current, [sourceFile]: nextPath }));
      onMetadataChanged();
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to update source path');
    } finally {
      setBusyKey(null);
    }
  };

  const handleSourceDelete = async (sourceFile: string) => {
    if (!window.confirm(`Delete metadata entry for ${sourceFile}?`)) {
      return;
    }

    setBusyKey(`source-delete:${sourceFile}`);
    setActionError(null);
    try {
      await deleteSourcePath(sourceFile);
      onMetadataChanged();
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to delete source metadata');
    } finally {
      setBusyKey(null);
    }
  };

  const handleTestSave = async (sourceFile: string, testIssue: MetadataTestIssue) => {
    const key = testIssueKey(sourceFile, testIssue);
    const nextPath = (testDrafts[key] ?? testIssue.testFile).trim();
    if (!nextPath) {
      setActionError('Replacement test path is required');
      return;
    }

    setBusyKey(`test-save:${key}`);
    setActionError(null);
    try {
      await updateTestPath(sourceFile, testIssue.testFile, testIssue.testName, nextPath);
      setTestDrafts((current) => ({ ...current, [key]: nextPath }));
      onMetadataChanged();
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to update test path');
    } finally {
      setBusyKey(null);
    }
  };

  const handleTestDelete = async (sourceFile: string, testIssue: MetadataTestIssue) => {
    if (!window.confirm(`Delete metadata for ${testIssue.testName} in ${testIssue.testFile}?`)) {
      return;
    }

    const key = testIssueKey(sourceFile, testIssue)
    setBusyKey(`test-delete:${key}`);
    setActionError(null);
    try {
      await deleteTestPath(sourceFile, testIssue.testFile, testIssue.testName);
      onMetadataChanged();
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to delete test metadata');
    } finally {
      setBusyKey(null);
    }
  };

  return (
    <div className="metadata-issues-panel">
      <div className="metadata-issues-summary">
        <strong>{issues.length}</strong> metadata entr{issues.length === 1 ? 'y has' : 'ies have'} invalid file references.
      </div>

      {actionError && (
        <div className="error-state">
          Error: {actionError}
        </div>
      )}

      <div className="metadata-issues-list">
        {issues.map((issue) => (
          <div key={issue.sourceFile} className="metadata-issue-card">
            <div className="metadata-issue-header">
              <div>
                <div className="metadata-issue-title">Source Entry</div>
                <div className="inline-code break-all">{issue.sourceFile}</div>
              </div>
              <div className="metadata-issue-badges">
                {!issue.sourceValid && (
                  <span className="badge badge-warning">invalid source</span>
                )}
                {issue.invalidTestIssues.length > 0 && (
                  <span className="badge badge-secondary">
                    {issue.invalidTestIssues.length} invalid test{issue.invalidTestIssues.length === 1 ? '' : 's'}
                  </span>
                )}
              </div>
            </div>

            {issue.sourceValid ? (
              <button
                className="inline-code-link metadata-link-button"
                onClick={() => onSourceFileClick(issue.sourceFile, 1)}
              >
                Open source file
              </button>
            ) : (
              <>
                <div className="metadata-issue-message">{issue.sourceMessage}</div>
                <div className="metadata-issue-edit-row">
                  <input
                    className="search-input"
                    value={sourceDrafts[issue.sourceFile] ?? issue.sourceFile}
                    onChange={(event) => setSourceDrafts((current) => ({
                      ...current,
                      [issue.sourceFile]: event.target.value,
                    }))}
                    placeholder="Enter replacement repo-relative source path"
                  />
                  <button
                    className="btn-control btn-control--active"
                    onClick={() => handleSourceSave(issue.sourceFile)}
                    disabled={busyKey === `source-save:${issue.sourceFile}`}
                  >
                    Save
                  </button>
                  <button
                    className="btn-control"
                    onClick={() => handleSourceDelete(issue.sourceFile)}
                    disabled={busyKey === `source-delete:${issue.sourceFile}`}
                  >
                    Delete entry
                  </button>
                </div>
              </>
            )}

            {(issue.commentsCount > 0 || issue.suggestionsCount > 0) && (
              <div className="metadata-issue-meta">
                {issue.commentsCount} comment{issue.commentsCount === 1 ? '' : 's'} and {issue.suggestionsCount} suggestion{issue.suggestionsCount === 1 ? '' : 's'} are attached to this source entry.
              </div>
            )}

            {issue.invalidTestIssues.length > 0 && (
              <div className="metadata-issue-tests">
                {issue.invalidTestIssues.map((testIssue) => {
                  const key = testIssueKey(issue.sourceFile, testIssue);
                  return (
                    <div key={key} className="metadata-test-issue">
                      <div className="metadata-test-issue-title">{testIssue.testName}</div>
                      <div className="inline-code break-all">{testIssue.testFile}</div>
                      <div className="metadata-issue-message">{testIssue.message}</div>
                      <div className="metadata-issue-edit-row">
                        <input
                          className="search-input"
                          value={testDrafts[key] ?? testIssue.testFile}
                          onChange={(event) => setTestDrafts((current) => ({
                            ...current,
                            [key]: event.target.value,
                          }))}
                          placeholder="Enter replacement repo-relative test path"
                        />
                        <button
                          className="btn-control btn-control--active"
                          onClick={() => handleTestSave(issue.sourceFile, testIssue)}
                          disabled={busyKey === `test-save:${key}`}
                        >
                          Save
                        </button>
                        <button
                          className="btn-control"
                          onClick={() => handleTestDelete(issue.sourceFile, testIssue)}
                          disabled={busyKey === `test-delete:${key}`}
                        >
                          Delete test
                        </button>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};
