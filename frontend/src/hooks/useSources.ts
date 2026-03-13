import { useState, useEffect } from 'react';
import { getSourceReferences } from '../api/client';
import type { SourceReference } from '../types';

export function useSources(path: string | null) {
  const [sources, setSources] = useState<SourceReference[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!path) {
      setSources([]);
      return;
    }

    setLoading(true);
    setError(null);

    getSourceReferences(path)
      .then(res => setSources(res.sources ?? []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [path]);

  return { sources, loading, error };
}
