import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useFiles } from './useFiles';
import * as client from '../api/client';
import type { ListFilesResponse } from '../types';

// Mock the API client
vi.mock('../api/client', () => ({
  listFiles: vi.fn(),
}));

describe('useFiles', () => {
  const mockListFiles = vi.mocked(client.listFiles);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should initialize with loading state', () => {
    mockListFiles.mockReturnValue(new Promise(() => {})); // Never resolves

    const { result } = renderHook(() => useFiles('.'));

    expect(result.current.loading).toBe(true);
    expect(result.current.files).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should fetch files successfully', async () => {
    const mockResponse: ListFilesResponse = {
      path: '.',
      files: [
        { name: 'file1.ts', path: 'file1.ts', isDir: false, modTime: '2024-01-01' },
        { name: 'dir1', path: 'dir1', isDir: true, modTime: '2024-01-01' },
      ],
    };

    mockListFiles.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useFiles('.'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.files).toEqual(mockResponse.files);
    expect(result.current.error).toBeNull();
    expect(mockListFiles).toHaveBeenCalledWith('.');
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch files';
    mockListFiles.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useFiles('.'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.files).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refetch when path changes', async () => {
    const mockResponse1: ListFilesResponse = {
      path: '.',
      files: [{ name: 'root.ts', path: 'root.ts', isDir: false, modTime: '2024-01-01' }],
    };

    const mockResponse2: ListFilesResponse = {
      path: 'src',
      files: [{ name: 'app.ts', path: 'src/app.ts', isDir: false, modTime: '2024-01-01' }],
    };

    mockListFiles
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result, rerender } = renderHook(({ path }) => useFiles(path), {
      initialProps: { path: '.' },
    });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.files).toEqual(mockResponse1.files);

    rerender({ path: 'src' });

    await waitFor(() => {
      expect(result.current.files).toEqual(mockResponse2.files);
    });

    expect(mockListFiles).toHaveBeenCalledTimes(2);
    expect(mockListFiles).toHaveBeenLastCalledWith('src');
  });

  it('should handle empty files array', async () => {
    const mockResponse: ListFilesResponse = {
      path: 'empty',
      files: [],
    };

    mockListFiles.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useFiles('empty'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.files).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should handle error with non-Error object', async () => {
    mockListFiles.mockRejectedValue('String error');

    const { result } = renderHook(() => useFiles('.'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle files with undefined response', async () => {
    mockListFiles.mockResolvedValue({ path: '.', files: undefined as unknown as [] });

    const { result } = renderHook(() => useFiles('.'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.files).toEqual([]);
  });
});
