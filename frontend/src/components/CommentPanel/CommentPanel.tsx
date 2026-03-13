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
    <div className="comment-panel-container">
      <div className="comment-panel-header">
        <h2 className="comment-panel-title">
          Comments
          {unresolvedCount > 0 && (
            <span className="badge badge-warning">
              {unresolvedCount}
            </span>
          )}
        </h2>
        <button
          onClick={onExportForAI}
          className="btn btn-primary"
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
        <div className="loading-state">
          Loading comments...
        </div>
      )}

      {error && (
        <div className="error-state">
          Error: {error}
        </div>
      )}

      {/* Add new comment form */}
      {selectedLine && (
        <form onSubmit={handleSubmit} className="comment-form">
          <div className="comment-form-wrapper">
            <div className="comment-form-label">
              Adding comment on line {selectedLine}
            </div>
            <textarea
              value={newCommentContent}
              onChange={(e) => setNewCommentContent(e.target.value)}
              placeholder="Enter your comment..."
              className="textarea"
            />
            <div className="comment-form-actions">
              <button
                type="submit"
                disabled={!newCommentContent.trim()}
                className={`btn btn-sm ${newCommentContent.trim() ? 'btn-primary' : 'btn-ghost'}`}
              >
                Add Comment
              </button>
            </div>
          </div>
        </form>
      )}

      {!selectedLine && comments.length === 0 && !loading && (
        <div className="empty-state">
          Click on a line number to add a comment
        </div>
      )}

      {/* Comments list */}
      <div className="flex-1 overflow-auto">
        {resolvedCount > 0 && (
          <label className="checkbox-label show-resolved-toggle">
            <input
              type="checkbox"
              checked={showResolved}
              onChange={(e) => setShowResolved(e.target.checked)}
            />
            Show resolved ({resolvedCount})
          </label>
        )}

        {sortedComments.length === 0 && !loading && comments.length > 0 && (
          <div className="empty-state">
            No unresolved comments
          </div>
        )}

        {sortedComments.map((comment) => (
          <div
            key={comment.id}
            className={`comment-item ${comment.resolved ? 'comment-item--resolved' : ''}`}
          >
            <div className="comment-item-header">
              <div className="comment-item-line">
                Line {comment.line}
                {comment.resolved && (
                  <span className="comment-item-resolved-badge">
                    Resolved
                  </span>
                )}
              </div>
              <div className="comment-item-actions">
                <button
                  onClick={() => onToggleResolved(comment.id)}
                  title={comment.resolved ? 'Mark as unresolved' : 'Mark as resolved'}
                  className={`btn-icon ${comment.resolved ? 'btn-icon--warning' : 'btn-icon--success'}`}
                >
                  {comment.resolved ? '↩' : '✓'}
                </button>
                {!comment.resolved && (
                  <button
                    onClick={() => startEditing(comment)}
                    title="Edit"
                    className="btn-icon"
                  >
                    ✎
                  </button>
                )}
                <button
                  onClick={() => onDeleteComment(comment.id)}
                  title="Delete"
                  className="btn-icon btn-icon--error"
                >
                  ×
                </button>
                <button
                  onClick={() => handleExportCommentForAI(comment)}
                  title={copiedCommentId === comment.id ? 'Copied!' : 'Copy for LLM agent'}
                  className={`btn-icon ${copiedCommentId === comment.id ? 'btn-icon--success' : ''}`}
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
                  className="textarea mb-sm"
                />
                <div className="comment-form-actions">
                  <button
                    onClick={cancelEdit}
                    className="btn btn-sm btn-ghost"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={saveEdit}
                    disabled={!editContent.trim()}
                    className={`btn btn-sm ${editContent.trim() ? 'btn-primary' : 'btn-ghost'}`}
                  >
                    Save
                  </button>
                </div>
              </div>
            ) : (
              <div className="comment-item-content">
                {comment.content}
              </div>
            )}

            <div className="comment-item-meta">
              {new Date(comment.createdAt).toLocaleString()}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
