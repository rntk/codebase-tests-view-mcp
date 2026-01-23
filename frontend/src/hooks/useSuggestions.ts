import { useState, useEffect } from 'react';
import { getSuggestions } from '../api/client';
import type { TestSuggestion } from '../types';

export function useSuggestions(path: string | null) {
  const [suggestions, setSuggestions] = useState<TestSuggestion[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!path) {
      setSuggestions([]);
      return;
    }

    setLoading(true);
    setError(null);

    getSuggestions(path)
      .then(res => setSuggestions(res.suggestions ?? []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [path]);

  return { suggestions, loading, error };
}
