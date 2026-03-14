import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useFileContent } from './useFileContent';
import * as client from '../api/client';
import type { FileResponse } from '../types';

vi.mock('../api/client', () => ({
  getFileContent: vi.fn(),
}));

describe('useFileContent', () => {
  const mockGetFileContent = vi.mocked(client.getFileContent);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not fetch when path is null', () => {
    const { result } = renderHook(() => useFileContent(null));

    expect(result.current.loading).toBe(false);
    expect(result.current.file).toBeNull();
    expect(result.current.error).toBeNull();
    expect(mockGetFileContent).not.toHaveBeenCalled();
  });

  it('should fetch file content when path is provided', async () => {
    const mockResponse: FileResponse = {
      file: {
        path: 'src/App.tsx',
        name: 'App.tsx',
        content: 'export default function App() {}',
        size: 32,
        modTime: '2024-01-01',
        mimeType: 'text/typescript',
      },
    };

    mockGetFileContent.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useFileContent('src/App.tsx'));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.file).toEqual(mockResponse.file);
    expect(result.current.error).toBeNull();
    expect(mockGetFileContent).toHaveBeenCalledWith('src/App.tsx');
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'File not found';
    mockGetFileContent.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useFileContent('missing.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.file).toBeNull();
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refetch when path changes', async () => {
    const mockResponse1: FileResponse = {
      file: {
        path: 'file1.ts',
        name: 'file1.ts',
        content: 'content1',
        size: 8,
        modTime: '2024-01-01',
        mimeType: 'text/typescript',
      },
    };

    const mockResponse2: FileResponse = {
      file: {
        path: 'file2.ts',
        name: 'file2.ts',
        content: 'content2',
        size: 8,
        modTime: '2024-01-01',
        mimeType: 'text/typescript',
      },
    };

    mockGetFileContent
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result, rerender } = renderHook(({ path }) => useFileContent(path), {
      initialProps: { path: 'file1.ts' },
    });

    await waitFor(() => {
      expect(result.current.file?.name).toBe('file1.ts');
    });

    rerender({ path: 'file2.ts' });

    await waitFor(() => {
      expect(result.current.file?.name).toBe('file2.ts');
    });

    expect(mockGetFileContent).toHaveBeenCalledTimes(2);
  });

  it('should clear file when path becomes null', async () => {
    const mockResponse: FileResponse = {
      file: {
        path: 'file.ts',
        name: 'file.ts',
        content: 'content',
        size: 7,
        modTime: '2024-01-01',
        mimeType: 'text/typescript',
      },
    };

    mockGetFileContent.mockResolvedValue(mockResponse);

    const { result, rerender } = renderHook<ReturnType<typeof useFileContent>, { path: string | null }>(({ path }) => useFileContent(path), {
      initialProps: { path: 'file.ts' },
    });

    await waitFor(() => {
      expect(result.current.file).not.toBeNull();
    });

    rerender({ path: null });

    expect(result.current.file).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  it('should handle error with non-Error object', async () => {
    mockGetFileContent.mockRejectedValue('Network error');

    const { result } = renderHook(() => useFileContent('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle file with metadata', async () => {
    const mockResponse: FileResponse = {
      file: {
        path: 'src/utils.ts',
        name: 'utils.ts',
        content: 'export function helper() {}',
        size: 27,
        modTime: '2024-01-01',
        mimeType: 'text/typescript',
        metadata: {
          tests: [],
        },
      },
    };

    mockGetFileContent.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useFileContent('src/utils.ts'));

    await waitFor(() => {
      expect(result.current.file).toEqual(mockResponse.file);
    });
  });
});
