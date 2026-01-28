import React, { useState } from 'react';
import type { Comment } from '../../types';

interface CommentPanelProps {
  comments: Comment[];
  loading: boolean;
  error: string | null;
  selectedLine?: number | null;
  filePath?: string | null;
  fileContent?: string | null;
  onAddComment: (line: number, content: string) => void;
  onUpdateComment: (commentId: string, content: string) => void;
  onDeleteComment: (commentId: string) => void;
  onToggleResolved: (commentId: string) => void;
  onExportForAI: () => void;
}

export const CommentPanel: React.FC<CommentPanelProps> = ({
  comments,
  loading,
  error,
  selectedLine,
  filePath,
  fileContent,
  onAddComment,
  onUpdateComment,
  onDeleteComment,
  onToggleResolved,
  onExportForAI,
}) => {
  const [newCommentContent, setNewCommentContent] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editContent, setEditContent] = useState('');
  const [showResolved, setShowResolved] = useState(false);
  const [copiedCommentId, setCopiedCommentId] = useState<string | null>(null);

  const filteredComments = showResolved 
    ? comments 
    : comments.filter(c => !c.resolved);

  const sortedComments = [...filteredComments].sort((a, b) => {
    // Sort by resolved status first, then by line number
    if (a.resolved !== b.resolved) {
      return a.resolved ? 1 : -1;
    }
    return a.line - b.line;
  });

  const unresolvedCount = comments.filter(c => !c.resolved).length;
  const resolvedCount = comments.filter(c => c.resolved).length;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedLine && newCommentContent.trim()) {
      onAddComment(selectedLine, newCommentContent.trim());
      setNewCommentContent('');
    }
  };

  const startEditing = (comment: Comment) => {
    setEditingId(comment.id);
    setEditContent(comment.content);
  };

  const saveEdit = () => {
    if (editingId && editContent.trim()) {
      onUpdateComment(editingId, editContent.trim());
      setEditingId(null);
      setEditContent('');
    }
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditContent('');
  };

  const handleExportCommentForAI = async (comment: Comment) => {
    if (!filePath || !fileContent) return;

    // Get the specific line of code being commented on
    const lines = fileContent.split('\n');
    const codeLine = lines[comment.line - 1] || '';

    const exportData = {
      file: filePath,
      line: comment.line,
      comment: comment.content,
      code: codeLine.trim(),
    };

    const formattedText = `File: ${exportData.file}
Line ${exportData.line}: ${exportData.code}

Comment: ${exportData.comment}

---
Please review this code and the associated comment. Fix any issues or make improvements as suggested.`;

    await navigator.clipboard.writeText(formattedText);
    setCopiedCommentId(comment.id);
    setTimeout(() => setCopiedCommentId(null), 2000);
  };

  return (
    <div className="comment-panel-container" style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: 'var(--space-md)',
        padding: '0 var(--space-sm)'
      }}>
        <h2 style={{
          margin: 0,
          fontSize: '16px',
          fontWeight: '600',
          color: 'var(--text-primary)'
        }}>
          Comments
          {unresolvedCount > 0 && (
            <span style={{
              marginLeft: '8px',
              padding: '2px 8px',
              backgroundColor: 'var(--warning)',
              color: 'white',
              borderRadius: '10px',
              fontSize: '11px',
            }}>
              {unresolvedCount}
            </span>
          )}
        </h2>
        <button
          onClick={onExportForAI}
          style={{
            padding: '6px 12px',
            backgroundColor: 'var(--accent-primary)',
            color: 'white',
            border: 'none',
            borderRadius: 'var(--radius-sm)',
            fontSize: '12px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '4px',
          }}
          title="Export comments and context for AI agent"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
            <polyline points="14,2 14,8 20,8" />
            <line x1="16" y1="13" x2="8" y2="13" />
            <line x1="16" y1="17" x2="8" y2="17" />
            <polyline points="10,9 9,9 8,9" />
          </svg>
          Export for AI
        </button>
      </div>

      {loading && (
        <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)' }}>
          Loading comments...
        </div>
      )}

      {error && (
        <div style={{
          padding: 'var(--space-md)',
          color: 'var(--error)',
          backgroundColor: '#fef2f2',
          borderRadius: 'var(--radius-md)',
          fontSize: '13px',
          marginBottom: 'var(--space-md)'
        }}>
          Error: {error}
        </div>
      )}

      {/* Add new comment form */}
      {selectedLine && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 'var(--space-md)' }}>
          <div style={{
            padding: 'var(--space-sm)',
            backgroundColor: 'var(--bg-secondary)',
            borderRadius: 'var(--radius-md)',
            border: '1px solid var(--border-color)',
          }}>
            <div style={{
              fontSize: '12px',
              color: 'var(--text-secondary)',
              marginBottom: 'var(--space-xs)',
            }}>
              Adding comment on line {selectedLine}
            </div>
            <textarea
              value={newCommentContent}
              onChange={(e) => setNewCommentContent(e.target.value)}
              placeholder="Enter your comment..."
              style={{
                width: '100%',
                minHeight: '60px',
                padding: 'var(--space-sm)',
                border: '1px solid var(--border-color)',
                borderRadius: 'var(--radius-sm)',
                fontSize: '13px',
                resize: 'vertical',
                fontFamily: 'inherit',
              }}
            />
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 'var(--space-sm)', marginTop: 'var(--space-sm)' }}>
              <button
                type="submit"
                disabled={!newCommentContent.trim()}
                style={{
                  padding: '6px 12px',
                  backgroundColor: newCommentContent.trim() ? 'var(--accent-primary)' : 'transparent',
                  color: newCommentContent.trim() ? 'white' : 'var(--text-secondary)',
                  border: newCommentContent.trim() ? 'none' : '1px solid var(--border-color)',
                  borderRadius: 'var(--radius-sm)',
                  fontSize: '12px',
                  cursor: newCommentContent.trim() ? 'pointer' : 'not-allowed',
                }}
              >
                Add Comment
              </button>
            </div>
          </div>
        </form>
      )}

      {!selectedLine && comments.length === 0 && !loading && (
        <div style={{
          padding: 'var(--space-lg) var(--space-md)',
          color: 'var(--text-tertiary)',
          textAlign: 'center',
          backgroundColor: 'var(--bg-primary)',
          borderRadius: 'var(--radius-md)',
          border: '1px dashed var(--border-color)',
        }}>
          Click on a line number to add a comment
        </div>
      )}

      {/* Comments list */}
      <div style={{ flex: 1, overflow: 'auto' }}>
        {resolvedCount > 0 && (
          <label style={{
            display: 'flex',
            alignItems: 'center',
            gap: 'var(--space-xs)',
            fontSize: '12px',
            color: 'var(--text-secondary)',
            marginBottom: 'var(--space-sm)',
            cursor: 'pointer',
          }}>
            <input
              type="checkbox"
              checked={showResolved}
              onChange={(e) => setShowResolved(e.target.checked)}
            />
            Show resolved ({resolvedCount})
          </label>
        )}

        {sortedComments.length === 0 && !loading && comments.length > 0 && (
          <div style={{
            padding: 'var(--space-md)',
            color: 'var(--text-tertiary)',
            textAlign: 'center',
            fontSize: '13px',
          }}>
            No unresolved comments
          </div>
        )}

        {sortedComments.map((comment) => (
          <div
            key={comment.id}
            style={{
              marginBottom: 'var(--space-md)',
              padding: 'var(--space-md)',
              backgroundColor: comment.resolved ? 'var(--bg-secondary)' : 'var(--bg-primary)',
              borderRadius: 'var(--radius-md)',
              border: '1px solid var(--border-color)',
              opacity: comment.resolved ? 0.7 : 1,
            }}
          >
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'flex-start',
              marginBottom: 'var(--space-xs)',
            }}>
              <div style={{
                fontSize: '12px',
                fontWeight: '600',
                color: 'var(--text-secondary)',
              }}>
                Line {comment.line}
                {comment.resolved && (
                  <span style={{
                    marginLeft: '8px',
                    padding: '2px 6px',
                    backgroundColor: 'var(--success)',
                    color: 'white',
                    borderRadius: '4px',
                    fontSize: '10px',
                  }}>
                    Resolved
                  </span>
                )}
              </div>
              <div style={{ display: 'flex', gap: 'var(--space-xs)' }}>
                <button
                  onClick={() => onToggleResolved(comment.id)}
                  title={comment.resolved ? 'Mark as unresolved' : 'Mark as resolved'}
                  style={{
                    padding: '4px',
                    backgroundColor: 'transparent',
                    border: 'none',
                    cursor: 'pointer',
                    color: comment.resolved ? 'var(--warning)' : 'var(--success)',
                  }}
                >
                  {comment.resolved ? '↩' : '✓'}
                </button>
                {!comment.resolved && (
                  <button
                    onClick={() => startEditing(comment)}
                    title="Edit"
                    style={{
                      padding: '4px',
                      backgroundColor: 'transparent',
                      border: 'none',
                      cursor: 'pointer',
                      color: 'var(--text-secondary)',
                    }}
                  >
                    ✎
                  </button>
                )}
                <button
                  onClick={() => onDeleteComment(comment.id)}
                  title="Delete"
                  style={{
                    padding: '4px',
                    backgroundColor: 'transparent',
                    border: 'none',
                    cursor: 'pointer',
                    color: 'var(--error)',
                  }}
                >
                  ×
                </button>
                <button
                  onClick={() => handleExportCommentForAI(comment)}
                  title={copiedCommentId === comment.id ? 'Copied!' : 'Copy for LLM agent'}
                  style={{
                    padding: '4px',
                    backgroundColor: 'transparent',
                    border: 'none',
                    cursor: 'pointer',
                    color: copiedCommentId === comment.id ? 'var(--success)' : 'var(--accent-primary)',
                  }}
                >
                  {copiedCommentId === comment.id ? (
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <polyline points="20,6 9,17 4,12" />
                    </svg>
                  ) : (
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                    </svg>
                  )}
                </button>
              </div>
            </div>

            {editingId === comment.id ? (
              <div>
                <textarea
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  style={{
                    width: '100%',
                    minHeight: '60px',
                    padding: 'var(--space-sm)',
                    border: '1px solid var(--border-color)',
                    borderRadius: 'var(--radius-sm)',
                    fontSize: '13px',
                    resize: 'vertical',
                    fontFamily: 'inherit',
                    marginBottom: 'var(--space-xs)',
                  }}
                />
                <div style={{ display: 'flex', gap: 'var(--space-sm)', justifyContent: 'flex-end' }}>
                  <button
                    onClick={cancelEdit}
                    style={{
                      padding: '4px 8px',
                      backgroundColor: 'transparent',
                      border: '1px solid var(--border-color)',
                      borderRadius: 'var(--radius-sm)',
                      fontSize: '12px',
                      cursor: 'pointer',
                    }}
                  >
                    Cancel
                  </button>
                  <button
                    onClick={saveEdit}
                    disabled={!editContent.trim()}
                    style={{
                      padding: '4px 8px',
                      backgroundColor: editContent.trim() ? 'var(--accent-primary)' : 'var(--bg-tertiary)',
                      color: editContent.trim() ? 'white' : 'var(--text-tertiary)',
                      border: 'none',
                      borderRadius: 'var(--radius-sm)',
                      fontSize: '12px',
                      cursor: editContent.trim() ? 'pointer' : 'not-allowed',
                    }}
                  >
                    Save
                  </button>
                </div>
              </div>
            ) : (
              <div style={{
                fontSize: '13px',
                color: 'var(--text-primary)',
                lineHeight: 1.5,
                whiteSpace: 'pre-wrap',
              }}>
                {comment.content}
              </div>
            )}

            <div style={{
              fontSize: '11px',
              color: 'var(--text-tertiary)',
              marginTop: 'var(--space-xs)',
            }}>
              {new Date(comment.createdAt).toLocaleString()}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
