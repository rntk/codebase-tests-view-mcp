import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useMetadataIssues } from './useMetadataIssues';
import * as client from '../api/client';
import type { MetadataIssuesResponse } from '../types';

vi.mock('../api/client', () => ({
  getMetadataIssues: vi.fn(),
}));

describe('useMetadataIssues', () => {
  const mockGetMetadataIssues = vi.mocked(client.getMetadataIssues);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should initialize with loading state', () => {
    mockGetMetadataIssues.mockReturnValue(new Promise(() => {}));

    const { result } = renderHook(() => useMetadataIssues());

    expect(result.current.loading).toBe(true);
    expect(result.current.issues).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should fetch metadata issues successfully', async () => {
    const mockResponse: MetadataIssuesResponse = {
      issues: [
        {
          sourceFile: 'src/old.ts',
          sourceValid: false,
          sourceIsAbsolute: false,
          sourceMessage: 'File not found',
          commentsCount: 0,
          invalidTestIssues: [],
        },
      ],
    };

    mockGetMetadataIssues.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.issues).toEqual(mockResponse.issues);
    expect(result.current.error).toBeNull();
    expect(mockGetMetadataIssues).toHaveBeenCalledTimes(1);
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch metadata issues';
    mockGetMetadataIssues.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.issues).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refresh when refresh function is called', async () => {
    const mockResponse1: MetadataIssuesResponse = {
      issues: [{ sourceFile: 'file1.ts', sourceValid: false, sourceIsAbsolute: false, sourceMessage: 'error1', commentsCount: 0, invalidTestIssues: [] }],
    };

    const mockResponse2: MetadataIssuesResponse = {
      issues: [
        { sourceFile: 'file1.ts', sourceValid: false, sourceIsAbsolute: false, sourceMessage: 'error1', commentsCount: 0, invalidTestIssues: [] },
        { sourceFile: 'file2.ts', sourceValid: false, sourceIsAbsolute: false, sourceMessage: 'error2', commentsCount: 0, invalidTestIssues: [] },
      ],
    };

    mockGetMetadataIssues
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.issues).toHaveLength(1);
    });

    act(() => {
      result.current.refresh();
    });

    await waitFor(() => {
      expect(result.current.issues).toHaveLength(2);
    });

    expect(mockGetMetadataIssues).toHaveBeenCalledTimes(2);
  });

  it('should handle error with non-Error object', async () => {
    mockGetMetadataIssues.mockRejectedValue('String error');

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle empty issues array', async () => {
    const mockResponse: MetadataIssuesResponse = {
      issues: [],
    };

    mockGetMetadataIssues.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.issues).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should handle issues with invalid test issues', async () => {
    const mockResponse: MetadataIssuesResponse = {
      issues: [
        {
          sourceFile: 'src/utils.ts',
          sourceValid: true,
          sourceIsAbsolute: false,
          commentsCount: 1,
          invalidTestIssues: [
            {
              testFile: 'src/utils.test.ts',
              testName: 'oldTest',
              isAbsolute: false,
              message: 'Test not found',
            },
          ],
        },
      ],
    };

    mockGetMetadataIssues.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.issues).toHaveLength(1);
    });

    const issue = result.current.issues[0];
    expect(issue.sourceValid).toBe(true);
    expect(issue.commentsCount).toBe(1);
    expect(issue.invalidTestIssues).toHaveLength(1);
    expect(issue.invalidTestIssues[0].testName).toBe('oldTest');
  });

  it('should handle issues with undefined response', async () => {
    mockGetMetadataIssues.mockResolvedValue({ issues: undefined as unknown as [] });

    const { result } = renderHook(() => useMetadataIssues());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.issues).toEqual([]);
  });
});
