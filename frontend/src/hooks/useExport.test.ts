import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useExport } from './useExport';
import * as client from '../api/client';
import type { ExportContextResponse } from '../types';

vi.mock('../api/client', () => ({
  exportContext: vi.fn(),
}));

describe('useExport', () => {
  const mockExportContext = vi.mocked(client.exportContext);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should initialize with default state', () => {
    const { result } = renderHook(() => useExport());

    expect(result.current.loading).toBe(false);
    expect(result.current.exportData).toBeNull();
    expect(result.current.error).toBeNull();
  });

  it('should perform export with default options', async () => {
    const mockResponse: ExportContextResponse = {
      sourceFile: 'src/App.tsx',
      codeContext: [
        {
          lineRange: { start: 1, end: 10 },
          code: 'const x = 1;',
          comments: [],
        },
      ],
      formatted: 'formatted output',
    };

    mockExportContext.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useExport());

    await act(async () => {
      await result.current.performExport('src/App.tsx');
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.exportData).toEqual(mockResponse);
    expect(result.current.error).toBeNull();
    expect(mockExportContext).toHaveBeenCalledWith('src/App.tsx', {
      includeTests: true,
      contextLines: 5,
    });
  });

  it('should perform export with custom options', async () => {
    const mockResponse: ExportContextResponse = {
      sourceFile: 'src/utils.ts',
      codeContext: [],
      tests: [],
      formatted: 'formatted',
    };

    mockExportContext.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useExport());

    await act(async () => {
      await result.current.performExport('src/utils.ts', {
        includeTests: false,
        contextLines: 10,
      });
    });

    expect(mockExportContext).toHaveBeenCalledWith('src/utils.ts', {
      includeTests: false,
      contextLines: 10,
    });
  });

  it('should handle export error', async () => {
    const errorMessage = 'Export failed';
    mockExportContext.mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useExport());

    await act(async () => {
      try {
        await result.current.performExport('src/App.tsx');
      } catch (e) {
        // Expected to throw
      }
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.exportData).toBeNull();
    expect(result.current.error).toBe(errorMessage);
  });

  it('should set loading state during export', async () => {
    mockExportContext.mockReturnValue(new Promise(() => {})); // Never resolves

    const { result } = renderHook(() => useExport());

    act(() => {
      result.current.performExport('src/App.tsx');
    });

    expect(result.current.loading).toBe(true);
  });

  it('should clear export data', async () => {
    const mockResponse: ExportContextResponse = {
      sourceFile: 'src/App.tsx',
      codeContext: [],
      formatted: 'formatted',
    };

    mockExportContext.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useExport());

    await act(async () => {
      await result.current.performExport('src/App.tsx');
    });

    expect(result.current.exportData).not.toBeNull();

    act(() => {
      result.current.clearExport();
    });

    expect(result.current.exportData).toBeNull();
    expect(result.current.error).toBeNull();
  });

  it('should handle error with non-Error object', async () => {
    mockExportContext.mockRejectedValue('String error');

    const { result } = renderHook(() => useExport());

    await act(async () => {
      try {
        await result.current.performExport('src/App.tsx');
      } catch (e) {
        // Expected to throw
      }
    });

    // When a string is thrown, the hook uses the fallback message
    expect(result.current.error).toBe('Failed to export context');
  });

  it('should handle export with tests', async () => {
    const mockResponse: ExportContextResponse = {
      sourceFile: 'src/utils.ts',
      codeContext: [],
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
      formatted: 'formatted with tests',
    };

    mockExportContext.mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useExport());

    await act(async () => {
      await result.current.performExport('src/utils.ts', {
        includeTests: true,
      });
    });

    expect(result.current.exportData?.tests).toHaveLength(1);
  });

  it('should merge options with defaults', async () => {
    mockExportContext.mockResolvedValue({
      sourceFile: 'src/App.tsx',
      codeContext: [],
      formatted: '',
    });

    const { result } = renderHook(() => useExport());

    await act(async () => {
      await result.current.performExport('src/App.tsx', { contextLines: 20 });
    });

    expect(mockExportContext).toHaveBeenCalledWith('src/App.tsx', {
      includeTests: true,
      contextLines: 20,
    });
  });
});
