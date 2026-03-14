import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useSources } from './useSources';
import * as client from '../api/client';
import type { TestFileResponse } from '../types';

vi.mock('../api/client', () => ({
  getSourceReferences: vi.fn(),
}));

describe('useSources', () => {
  const mockGetSourceReferences = vi.mocked(client.getSourceReferences);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not fetch when path is null', () => {
    const { result } = renderHook(() => useSources(null));

    expect(result.current.loading).toBe(false);
    expect(result.current.sources).toEqual([]);
    expect(result.current.error).toBeNull();
    expect(mockGetSourceReferences).not.toHaveBeenCalled();
  });

  it('should fetch source references when path is provided', async () => {
    const mockResponse: TestFileResponse = {
      testFile: 'src/utils.test.ts',
      sources: [
        {
          sourceFile: 'src/utils.ts',
          functionName: 'helper',
          coveredLines: { start: 5, end: 10 },
          testName: 'should work',
          lineRange: { start: 1, end: 3 },
        },
      ],
    };

    mockGetSourceReferences.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useSources('src/utils.test.ts'));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.sources).toEqual(mockResponse.sources);
    expect(result.current.error).toBeNull();
    expect(mockGetSourceReferences).toHaveBeenCalledWith('src/utils.test.ts');
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch sources';
    mockGetSourceReferences.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useSources('test.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.sources).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refetch when path changes', async () => {
    const mockResponse1: TestFileResponse = {
      testFile: 'test1.test.ts',
      sources: [{ sourceFile: 'src1.ts', functionName: 'fn1', coveredLines: { start: 1, end: 5 }, testName: 'test1', lineRange: { start: 1, end: 1 } }],
    };

    const mockResponse2: TestFileResponse = {
      testFile: 'test2.test.ts',
      sources: [{ sourceFile: 'src2.ts', functionName: 'fn2', coveredLines: { start: 10, end: 15 }, testName: 'test2', lineRange: { start: 1, end: 1 } }],
    };

    mockGetSourceReferences
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result, rerender } = renderHook(({ path }) => useSources(path), {
      initialProps: { path: 'test1.test.ts' },
    });

    await waitFor(() => {
      expect(result.current.sources).toHaveLength(1);
      expect(result.current.sources[0].functionName).toBe('fn1');
    });

    rerender({ path: 'test2.test.ts' });

    await waitFor(() => {
      expect(result.current.sources).toHaveLength(1);
      expect(result.current.sources[0].functionName).toBe('fn2');
    });

    expect(mockGetSourceReferences).toHaveBeenCalledTimes(2);
  });

  it('should clear sources when path becomes null', async () => {
    const mockResponse: TestFileResponse = {
      testFile: 'test.ts',
      sources: [{ sourceFile: 'src.ts', functionName: 'fn', coveredLines: { start: 1, end: 5 }, testName: 'test', lineRange: { start: 1, end: 1 } }],
    };

    mockGetSourceReferences.mockResolvedValue(mockResponse);

    const { result, rerender } = renderHook<ReturnType<typeof useSources>, { path: string | null }>(({ path }) => useSources(path), {
      initialProps: { path: 'test.ts' },
    });

    await waitFor(() => {
      expect(result.current.sources).toHaveLength(1);
    });

    rerender({ path: null });

    expect(result.current.sources).toEqual([]);
    expect(result.current.loading).toBe(false);
  });

  it('should handle empty sources array', async () => {
    const mockResponse: TestFileResponse = {
      testFile: 'test.ts',
      sources: [],
    };

    mockGetSourceReferences.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useSources('test.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.sources).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should handle sources with undefined response', async () => {
    mockGetSourceReferences.mockResolvedValue({ testFile: 'test.ts', sources: undefined as unknown as [] });

    const { result } = renderHook(() => useSources('test.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.sources).toEqual([]);
  });

  it('should handle error with non-Error object', async () => {
    mockGetSourceReferences.mockRejectedValue('String error');

    const { result } = renderHook(() => useSources('test.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle sources with full details', async () => {
    const mockResponse: TestFileResponse = {
      testFile: 'src/utils.test.ts',
      sources: [
        {
          sourceFile: 'src/utils.ts',
          functionName: 'formatDate',
          coveredLines: { start: 5, end: 15 },
          testName: 'should format date',
          lineRange: { start: 1, end: 10 },
          comment: 'Tests date formatting',
          inputLines: { start: 2, end: 4 },
          outputLines: { start: 7, end: 9 },
        },
      ],
    };

    mockGetSourceReferences.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useSources('src/utils.test.ts'));

    await waitFor(() => {
      expect(result.current.sources).toEqual(mockResponse.sources);
    });

    const source = result.current.sources[0];
    expect(source.comment).toBe('Tests date formatting');
    expect(source.inputLines).toEqual({ start: 2, end: 4 });
    expect(source.outputLines).toEqual({ start: 7, end: 9 });
  });
});
