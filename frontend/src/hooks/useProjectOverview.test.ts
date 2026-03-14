import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useProjectOverview } from './useProjectOverview';
import * as client from '../api/client';
import type { OverviewResponse } from '../types';

vi.mock('../api/client', () => ({
  getProjectOverview: vi.fn(),
}));

describe('useProjectOverview', () => {
  const mockGetProjectOverview = vi.mocked(client.getProjectOverview);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should initialize with loading state', () => {
    mockGetProjectOverview.mockReturnValue(new Promise(() => {}));

    const { result } = renderHook(() => useProjectOverview());

    expect(result.current.loading).toBe(true);
    expect(result.current.overview).toBeNull();
    expect(result.current.error).toBeNull();
  });

  it('should fetch overview successfully', async () => {
    const mockResponse: OverviewResponse = {
      totalTests: 100,
      totalFunctions: 50,
      totalSourceFiles: 25,
      totalTestFiles: 20,
      functions: [
        {
          functionName: 'helper',
          sourceFile: 'src/utils.ts',
          testCount: 3,
          tests: [],
        },
      ],
      testsBySourceFile: {
        'src/utils.ts': [],
      },
    };

    mockGetProjectOverview.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useProjectOverview());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.overview).toEqual(mockResponse);
    expect(result.current.error).toBeNull();
    expect(mockGetProjectOverview).toHaveBeenCalledTimes(1);
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch overview';
    mockGetProjectOverview.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useProjectOverview());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.overview).toBeNull();
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refresh when refresh function is called', async () => {
    const mockResponse1: OverviewResponse = {
      totalTests: 100,
      totalFunctions: 50,
      totalSourceFiles: 25,
      totalTestFiles: 20,
      functions: [],
      testsBySourceFile: {},
    };

    const mockResponse2: OverviewResponse = {
      totalTests: 150,
      totalFunctions: 60,
      totalSourceFiles: 30,
      totalTestFiles: 25,
      functions: [],
      testsBySourceFile: {},
    };

    mockGetProjectOverview
      .mockResolvedValueOnce(mockResponse1)
      .mockResolvedValueOnce(mockResponse2);

    const { result } = renderHook(() => useProjectOverview());

    await waitFor(() => {
      expect(result.current.overview?.totalTests).toBe(100);
    });

    act(() => {
      result.current.refresh();
    });

    await waitFor(() => {
      expect(result.current.overview?.totalTests).toBe(150);
    });

    expect(mockGetProjectOverview).toHaveBeenCalledTimes(2);
  });

  it('should handle error with non-Error object', async () => {
    mockGetProjectOverview.mockRejectedValue('String error');

    const { result } = renderHook(() => useProjectOverview());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // When a string is thrown, err.message is undefined
    expect(result.current.error).toBeUndefined();
  });

  it('should handle overview with empty arrays', async () => {
    const mockResponse: OverviewResponse = {
      totalTests: 0,
      totalFunctions: 0,
      totalSourceFiles: 0,
      totalTestFiles: 0,
      functions: [],
      testsBySourceFile: {},
    };

    mockGetProjectOverview.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useProjectOverview());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.overview).toEqual(mockResponse);
  });

  it('should handle overview with multiple functions', async () => {
    const mockResponse: OverviewResponse = {
      totalTests: 10,
      totalFunctions: 3,
      totalSourceFiles: 2,
      totalTestFiles: 2,
      functions: [
        { functionName: 'fn1', sourceFile: 'file1.ts', testCount: 4, tests: [] },
        { functionName: 'fn2', sourceFile: 'file1.ts', testCount: 3, tests: [] },
        { functionName: 'fn3', sourceFile: 'file2.ts', testCount: 3, tests: [] },
      ],
      testsBySourceFile: {
        'file1.ts': [],
        'file2.ts': [],
      },
    };

    mockGetProjectOverview.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useProjectOverview());

    await waitFor(() => {
      expect(result.current.overview?.functions).toHaveLength(3);
    });

    expect(result.current.overview?.totalFunctions).toBe(3);
    expect(result.current.overview?.totalTests).toBe(10);
  });
});
