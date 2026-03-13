import { useState, useEffect } from 'react';
import { getProjectOverview } from '../api/client';
import type { OverviewResponse } from '../types';

export function useProjectOverview() {
  const [overview, setOverview] = useState<OverviewResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    setLoading(true);
    setError(null);

    getProjectOverview()
      .then(res => setOverview(res))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [refreshKey]);

  return {
    overview,
    loading,
    error,
    refresh: () => setRefreshKey(key => key + 1),
  };
}
