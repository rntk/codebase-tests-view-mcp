import { useState, useEffect, useCallback } from 'react';
import type { Comment, CommentRequest } from '../types';
import {
  getComments,
  createComment,
  updateComment,
  deleteComment,
  toggleCommentResolved,
} from '../api/client';

interface UseCommentsResult {
  comments: Comment[];
  loading: boolean;
  error: string | null;
  addComment: (line: number, content: string) => Promise<void>;
  editComment: (commentId: string, content: string) => Promise<void>;
  removeComment: (commentId: string) => Promise<void>;
  toggleResolved: (commentId: string) => Promise<void>;
  refresh: () => void;
}

export function useComments(filePath: string | null): UseCommentsResult {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const refresh = useCallback(() => {
    setRefreshKey(prev => prev + 1);
  }, []);

  useEffect(() => {
    if (!filePath) {
      setComments([]);
      return;
    }

    let cancelled = false;

    async function fetchComments() {
      setLoading(true);
      setError(null);

      try {
        const response = await getComments(filePath!);
        if (!cancelled) {
          setComments(response.comments);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to fetch comments');
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    fetchComments();

    return () => {
      cancelled = true;
    };
  }, [filePath, refreshKey]);

  const addComment = useCallback(async (line: number, content: string) => {
    if (!filePath) return;

    try {
      const request: CommentRequest = {
        line,
        content,
        contextLines: { start: Math.max(1, line - 3), end: line + 3 },
      };
      await createComment(filePath, request);
      refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create comment');
      throw err;
    }
  }, [filePath, refresh]);

  const editComment = useCallback(async (commentId: string, content: string) => {
    if (!filePath) return;

    try {
      await updateComment(filePath, commentId, content);
      refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update comment');
      throw err;
    }
  }, [filePath, refresh]);

  const removeComment = useCallback(async (commentId: string) => {
    if (!filePath) return;

    try {
      await deleteComment(filePath, commentId);
      refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete comment');
      throw err;
    }
  }, [filePath, refresh]);

  const toggleResolved = useCallback(async (commentId: string) => {
    if (!filePath) return;

    try {
      await toggleCommentResolved(filePath, commentId);
      refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to toggle resolved status');
      throw err;
    }
  }, [filePath, refresh]);

  return {
    comments,
    loading,
    error,
    addComment,
    editComment,
    removeComment,
    toggleResolved,
    refresh,
  };
}
