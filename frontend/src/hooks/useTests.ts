import { useState, useEffect } from 'react';
import { getRelatedTests } from '../api/client';
import type { TestDetail } from '../types';

export function useTests(path: string | null) {
  const [tests, setTests] = useState<TestDetail[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!path) {
      setTests([]);
      return;
    }

    setLoading(true);
    setError(null);

    getRelatedTests(path)
      .then(res => setTests(res.tests))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [path]);

  return { tests, loading, error };
}
