import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useTests } from './useTests';
import * as client from '../api/client';
import type { TestsResponse } from '../types';

vi.mock('../api/client', () => ({
  getRelatedTests: vi.fn(),
}));

describe('useTests', () => {
  const mockGetRelatedTests = vi.mocked(client.getRelatedTests);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not fetch when path is null', () => {
    const { result } = renderHook(() => useTests(null));

    expect(result.current.loading).toBe(false);
    expect(result.current.tests).toEqual([]);
    expect(result.current.error).toBeNull();
    expect(mockGetRelatedTests).not.toHaveBeenCalled();
  });

  it('should fetch tests when path is provided', async () => {
    const mockResponse: TestsResponse = {
      sourceFile: 'src/utils.ts',
      tests: [
        {
          functionName: 'helper',
          testFile: 'src/utils.test.ts',
          testName: 'should work',
          content: 'test("should work", () => {})',
          lineRange: { start: 1, end: 3 },
          coveredLines: { start: 5, end: 10 },
        },
      ],
    };

    mockGetRelatedTests.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useTests('src/utils.ts'));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.tests).toEqual(mockResponse.tests);
    expect(result.current.error).toBeNull();
    expect(mockGetRelatedTests).toHaveBeenCalledWith('src/utils.ts');
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch tests';
    mockGetRelatedTests.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useTests('src/utils.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.tests).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refetch when path changes', async () => {
    const mockResponse1: TestsResponse = {
      sourceFile: 'file1.ts',
      tests: [{ functionName: 'fn1', testFile: 'file1.test.ts', testName: 'test1', content: '', lineRange: { start: 1, end: 1 }, coveredLines: { start: 1, end: 1 } }],
    };

    const mockResponse2: TestsResponse = {
      sourceFile: 'file2.ts',
      tests: [{ functionName: 'fn2', testFile: 'file2.test.ts', testName: 'test2', content: '', lineRange: { start: 1, end: 1 }, coveredLines: { start: 1, end: 1 } }],
    };

    mockGetRelatedTests
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result, rerender } = renderHook(({ path }) => useTests(path), {
      initialProps: { path: 'file1.ts' },
    });

    await waitFor(() => {
      expect(result.current.tests).toHaveLength(1);
      expect(result.current.tests[0].functionName).toBe('fn1');
    });

    rerender({ path: 'file2.ts' });

    await waitFor(() => {
      expect(result.current.tests).toHaveLength(1);
      expect(result.current.tests[0].functionName).toBe('fn2');
    });

    expect(mockGetRelatedTests).toHaveBeenCalledTimes(2);
  });

  it('should clear tests when path becomes null', async () => {
    const mockResponse: TestsResponse = {
      sourceFile: 'file.ts',
      tests: [{ functionName: 'fn', testFile: 'file.test.ts', testName: 'test', content: '', lineRange: { start: 1, end: 1 }, coveredLines: { start: 1, end: 1 } }],
    };

    mockGetRelatedTests.mockResolvedValue(mockResponse);

    const { result, rerender } = renderHook<ReturnType<typeof useTests>, { path: string | null }>(({ path }) => useTests(path), {
      initialProps: { path: 'file.ts' },
    });

    await waitFor(() => {
      expect(result.current.tests).toHaveLength(1);
    });

    rerender({ path: null });

    expect(result.current.tests).toEqual([]);
    expect(result.current.loading).toBe(false);
  });

  it('should handle empty tests array', async () => {
    const mockResponse: TestsResponse = {
      sourceFile: 'file.ts',
      tests: [],
    };

    mockGetRelatedTests.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useTests('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.tests).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should handle tests with undefined response', async () => {
    mockGetRelatedTests.mockResolvedValue({ sourceFile: 'file.ts', tests: undefined as unknown as [] });

    const { result } = renderHook(() => useTests('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.tests).toEqual([]);
  });

  it('should handle error with non-Error object', async () => {
    mockGetRelatedTests.mockRejectedValue('String error');

    const { result } = renderHook(() => useTests('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle tests with full details', async () => {
    const mockResponse: TestsResponse = {
      sourceFile: 'src/utils.ts',
      tests: [
        {
          functionName: 'formatDate',
          testFile: 'src/utils.test.ts',
          testName: 'should format date correctly',
          content: 'test("should format date correctly", () => {})',
          lineRange: { start: 1, end: 10 },
          coveredLines: { start: 5, end: 15 },
          comment: 'Tests date formatting',
          inputData: '2024-01-01',
          inputLines: { start: 2, end: 4 },
          expectedOutput: 'Jan 1, 2024',
          outputLines: { start: 7, end: 9 },
        },
      ],
    };

    mockGetRelatedTests.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useTests('src/utils.ts'));

    await waitFor(() => {
      expect(result.current.tests).toEqual(mockResponse.tests);
    });

    const test = result.current.tests[0];
    expect(test.comment).toBe('Tests date formatting');
    expect(test.inputData).toBe('2024-01-01');
    expect(test.expectedOutput).toBe('Jan 1, 2024');
  });
});
