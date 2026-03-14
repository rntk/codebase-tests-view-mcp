import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useSuggestions } from './useSuggestions';
import * as client from '../api/client';
import type { SuggestionsResponse } from '../types';

vi.mock('../api/client', () => ({
  getSuggestions: vi.fn(),
}));

describe('useSuggestions', () => {
  const mockGetSuggestions = vi.mocked(client.getSuggestions);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not fetch when path is null', () => {
    const { result } = renderHook(() => useSuggestions(null));

    expect(result.current.loading).toBe(false);
    expect(result.current.suggestions).toEqual([]);
    expect(result.current.error).toBeNull();
    expect(mockGetSuggestions).not.toHaveBeenCalled();
  });

  it('should fetch suggestions when path is provided', async () => {
    const mockResponse: SuggestionsResponse = {
      sourceFile: 'src/utils.ts',
      suggestions: [
        {
          sourceFile: 'src/utils.ts',
          functionName: 'untestedFunction',
          targetLines: { start: 10, end: 20 },
          reason: 'Function lacks test coverage',
          suggestedName: 'should handle untestedFunction',
          testSkeleton: 'test("should handle untestedFunction", () => {})',
          priority: 'high',
        },
      ],
    };

    mockGetSuggestions.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useSuggestions('src/utils.ts'));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.suggestions).toEqual(mockResponse.suggestions);
    expect(result.current.error).toBeNull();
    expect(mockGetSuggestions).toHaveBeenCalledWith('src/utils.ts');
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch suggestions';
    mockGetSuggestions.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useSuggestions('src/utils.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.suggestions).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refetch when path changes', async () => {
    const mockResponse1: SuggestionsResponse = {
      sourceFile: 'file1.ts',
      suggestions: [{ sourceFile: 'file1.ts', targetLines: { start: 1, end: 5 }, reason: 'reason1', suggestedName: 'test1', testSkeleton: '', priority: 'medium' }],
    };

    const mockResponse2: SuggestionsResponse = {
      sourceFile: 'file2.ts',
      suggestions: [{ sourceFile: 'file2.ts', targetLines: { start: 10, end: 15 }, reason: 'reason2', suggestedName: 'test2', testSkeleton: '', priority: 'high' }],
    };

    mockGetSuggestions
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result, rerender } = renderHook(({ path }) => useSuggestions(path), {
      initialProps: { path: 'file1.ts' },
    });

    await waitFor(() => {
      expect(result.current.suggestions).toHaveLength(1);
      expect(result.current.suggestions[0].suggestedName).toBe('test1');
    });

    rerender({ path: 'file2.ts' });

    await waitFor(() => {
      expect(result.current.suggestions).toHaveLength(1);
      expect(result.current.suggestions[0].suggestedName).toBe('test2');
    });

    expect(mockGetSuggestions).toHaveBeenCalledTimes(2);
  });

  it('should clear suggestions when path becomes null', async () => {
    const mockResponse: SuggestionsResponse = {
      sourceFile: 'file.ts',
      suggestions: [{ sourceFile: 'file.ts', targetLines: { start: 1, end: 5 }, reason: 'reason', suggestedName: 'test', testSkeleton: '', priority: 'low' }],
    };

    mockGetSuggestions.mockResolvedValue(mockResponse);

    const { result, rerender } = renderHook<ReturnType<typeof useSuggestions>, { path: string | null }>(({ path }) => useSuggestions(path), {
      initialProps: { path: 'file.ts' },
    });

    await waitFor(() => {
      expect(result.current.suggestions).toHaveLength(1);
    });

    rerender({ path: null });

    expect(result.current.suggestions).toEqual([]);
    expect(result.current.loading).toBe(false);
  });

  it('should handle empty suggestions array', async () => {
    const mockResponse: SuggestionsResponse = {
      sourceFile: 'file.ts',
      suggestions: [],
    };

    mockGetSuggestions.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useSuggestions('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.suggestions).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should handle suggestions with undefined response', async () => {
    mockGetSuggestions.mockResolvedValue({ sourceFile: 'file.ts', suggestions: undefined as unknown as [] });

    const { result } = renderHook(() => useSuggestions('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.suggestions).toEqual([]);
  });

  it('should handle error with non-Error object', async () => {
    mockGetSuggestions.mockRejectedValue('String error');

    const { result } = renderHook(() => useSuggestions('file.ts'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle all priority levels', async () => {
    const mockResponse: SuggestionsResponse = {
      sourceFile: 'src/utils.ts',
      suggestions: [
        { sourceFile: 'src/utils.ts', targetLines: { start: 1, end: 5 }, reason: 'r1', suggestedName: 'high', testSkeleton: '', priority: 'high' },
        { sourceFile: 'src/utils.ts', targetLines: { start: 6, end: 10 }, reason: 'r2', suggestedName: 'medium', testSkeleton: '', priority: 'medium' },
        { sourceFile: 'src/utils.ts', targetLines: { start: 11, end: 15 }, reason: 'r3', suggestedName: 'low', testSkeleton: '', priority: 'low' },
      ],
    };

    mockGetSuggestions.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useSuggestions('src/utils.ts'));

    await waitFor(() => {
      expect(result.current.suggestions).toHaveLength(3);
    });

    expect(result.current.suggestions[0].priority).toBe('high');
    expect(result.current.suggestions[1].priority).toBe('medium');
    expect(result.current.suggestions[2].priority).toBe('low');
  });
});
