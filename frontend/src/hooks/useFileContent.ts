import { useState, useEffect } from 'react';
import { getFileContent } from '../api/client';
import type { FileContent } from '../types';

export function useFileContent(path: string | null) {
  const [file, setFile] = useState<FileContent | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!path) {
      setFile(null);
      return;
    }

    setLoading(true);
    setError(null);

    getFileContent(path)
      .then(res => setFile(res.file))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [path]);

  return { file, loading, error };
}
