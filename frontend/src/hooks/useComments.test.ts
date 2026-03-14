import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useComments } from './useComments';
import * as client from '../api/client';
import type { CommentsResponse, Comment } from '../types';

vi.mock('../api/client', () => ({
  getComments: vi.fn(),
  createComment: vi.fn(),
  updateComment: vi.fn(),
  deleteComment: vi.fn(),
  toggleCommentResolved: vi.fn(),
}));

describe('useComments', () => {
  const mockGetComments = vi.mocked(client.getComments);
  const mockCreateComment = vi.mocked(client.createComment);
  const mockUpdateComment = vi.mocked(client.updateComment);
  const mockDeleteComment = vi.mocked(client.deleteComment);
  const mockToggleCommentResolved = vi.mocked(client.toggleCommentResolved);

  const mockComment: Comment = {
    id: 'comment-1',
    line: 10,
    content: 'Test comment',
    author: 'user',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    resolved: false,
    contextLines: { start: 7, end: 13 },
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not fetch when filePath is null', () => {
    const { result } = renderHook(() => useComments(null));

    expect(result.current.loading).toBe(false);
    expect(result.current.comments).toEqual([]);
    expect(result.current.error).toBeNull();
    expect(mockGetComments).not.toHaveBeenCalled();
  });

  it('should fetch comments when filePath is provided', async () => {
    const mockResponse: CommentsResponse = {
      sourceFile: 'src/App.tsx',
      comments: [mockComment],
    };

    mockGetComments.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useComments('src/App.tsx'));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.comments).toEqual([mockComment]);
    expect(result.current.error).toBeNull();
    expect(mockGetComments).toHaveBeenCalledWith('src/App.tsx');
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch comments';
    mockGetComments.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.comments).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should add a comment', async () => {
    const mockResponse: CommentsResponse = {
      sourceFile: 'src/App.tsx',
      comments: [mockComment],
    };

    mockGetComments.mockResolvedValue(mockResponse);
    mockCreateComment.mockResolvedValue({ comment: mockComment });

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.addComment(10, 'New comment');
    });

    expect(mockCreateComment).toHaveBeenCalledWith('src/App.tsx', {
      line: 10,
      content: 'New comment',
      contextLines: { start: 7, end: 13 },
    });
  });

  it('should not add comment when filePath is null', async () => {
    const { result } = renderHook(() => useComments(null));

    await act(async () => {
      await result.current.addComment(10, 'New comment');
    });

    expect(mockCreateComment).not.toHaveBeenCalled();
  });

  it('should handle add comment error', async () => {
    const mockResponse: CommentsResponse = {
      sourceFile: 'src/App.tsx',
      comments: [],
    };

    mockGetComments.mockResolvedValue(mockResponse);
    mockCreateComment.mockRejectedValue(new Error('Failed to create'));

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      try {
        await result.current.addComment(10, 'New comment');
      } catch (e) {
        // Expected to throw
      }
    });

    expect(result.current.error).toBe('Failed to create');
  });

  it('should edit a comment', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });
    mockUpdateComment.mockResolvedValue(undefined);

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.editComment('comment-1', 'Updated content');
    });

    expect(mockUpdateComment).toHaveBeenCalledWith('src/App.tsx', 'comment-1', 'Updated content');
  });

  it('should not edit comment when filePath is null', async () => {
    const { result } = renderHook(() => useComments(null));

    await act(async () => {
      await result.current.editComment('comment-1', 'Updated');
    });

    expect(mockUpdateComment).not.toHaveBeenCalled();
  });

  it('should handle edit comment error', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });
    mockUpdateComment.mockRejectedValue(new Error('Failed to update'));

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      try {
        await result.current.editComment('comment-1', 'Updated');
      } catch (e) {
        // Expected to throw
      }
    });

    expect(result.current.error).toBe('Failed to update');
  });

  it('should remove a comment', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });
    mockDeleteComment.mockResolvedValue(undefined);

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.removeComment('comment-1');
    });

    expect(mockDeleteComment).toHaveBeenCalledWith('src/App.tsx', 'comment-1');
  });

  it('should not remove comment when filePath is null', async () => {
    const { result } = renderHook(() => useComments(null));

    await act(async () => {
      await result.current.removeComment('comment-1');
    });

    expect(mockDeleteComment).not.toHaveBeenCalled();
  });

  it('should handle remove comment error', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });
    mockDeleteComment.mockRejectedValue(new Error('Failed to delete'));

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      try {
        await result.current.removeComment('comment-1');
      } catch (e) {
        // Expected to throw
      }
    });

    expect(result.current.error).toBe('Failed to delete');
  });

  it('should toggle resolved status', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });
    mockToggleCommentResolved.mockResolvedValue(undefined);

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.toggleResolved('comment-1');
    });

    expect(mockToggleCommentResolved).toHaveBeenCalledWith('src/App.tsx', 'comment-1');
  });

  it('should not toggle resolved when filePath is null', async () => {
    const { result } = renderHook(() => useComments(null));

    await act(async () => {
      await result.current.toggleResolved('comment-1');
    });

    expect(mockToggleCommentResolved).not.toHaveBeenCalled();
  });

  it('should handle toggle resolved error', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });
    mockToggleCommentResolved.mockRejectedValue(new Error('Failed to toggle'));

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      try {
        await result.current.toggleResolved('comment-1');
      } catch (e) {
        // Expected to throw
      }
    });

    expect(result.current.error).toBe('Failed to toggle');
  });

  it('should refresh comments', async () => {
    const mockResponse1: CommentsResponse = {
      sourceFile: 'src/App.tsx',
      comments: [mockComment],
    };

    const mockResponse2: CommentsResponse = {
      sourceFile: 'src/App.tsx',
      comments: [mockComment, { ...mockComment, id: 'comment-2', content: 'Second comment' }],
    };

    mockGetComments
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.comments).toHaveLength(1);
    });

    act(() => {
      result.current.refresh();
    });

    await waitFor(() => {
      expect(result.current.comments).toHaveLength(2);
    });

    expect(mockGetComments).toHaveBeenCalledTimes(2);
  });

  it('should clear comments when filePath becomes null', async () => {
    mockGetComments.mockResolvedValue({ sourceFile: 'src/App.tsx', comments: [mockComment] });

    const { result, rerender } = renderHook<ReturnType<typeof useComments>, { path: string | null }>(({ path }) => useComments(path), {
      initialProps: { path: 'src/App.tsx' },
    });

    await waitFor(() => {
      expect(result.current.comments).toHaveLength(1);
    });

    rerender({ path: null });

    expect(result.current.comments).toEqual([]);
    expect(result.current.loading).toBe(false);
  });

  it('should handle non-Error exceptions', async () => {
    mockGetComments.mockRejectedValue('String error');

    const { result } = renderHook(() => useComments('src/App.tsx'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Failed to fetch comments');
  });
});
