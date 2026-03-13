import { useEffect, useState } from 'react';
import { getMetadataIssues } from '../api/client';
import type { MetadataIssue } from '../types';

export function useMetadataIssues() {
  const [issues, setIssues] = useState<MetadataIssue[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    setLoading(true);
    setError(null);

    getMetadataIssues()
      .then((response) => setIssues(response.issues ?? []))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [refreshKey]);

  return {
    issues,
    loading,
    error,
    refresh: () => setRefreshKey((key) => key + 1),
  };
}
