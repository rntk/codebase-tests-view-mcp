import { useState, useCallback } from 'react';
import type { ExportContextRequest, ExportContextResponse } from '../types';
import { exportContext } from '../api/client';

interface UseExportResult {
  exportData: ExportContextResponse | null;
  loading: boolean;
  error: string | null;
  performExport: (filePath: string, options?: Partial<ExportContextRequest>) => Promise<void>;
  clearExport: () => void;
}

export function useExport(): UseExportResult {
  const [exportData, setExportData] = useState<ExportContextResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const performExport = useCallback(async (
    filePath: string, 
    options: Partial<ExportContextRequest> = {}
  ) => {
    setLoading(true);
    setError(null);

    try {
      const request: ExportContextRequest = {
        includeTests: true,
        includeSuggestions: true,
        contextLines: 5,
        ...options,
      };

      const response = await exportContext(filePath, request);
      setExportData(response);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to export context';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const clearExport = useCallback(() => {
    setExportData(null);
    setError(null);
  }, []);

  return {
    exportData,
    loading,
    error,
    performExport,
    clearExport,
  };
}
